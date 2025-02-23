package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmad-abuziad/transaction-processor/internal/data"
	"github.com/ahmad-abuziad/transaction-processor/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Status": "OK",
	})
}

func (app *application) newSalesTransaction(c *gin.Context) {
	// parse
	var input struct {
		TenantID     int64     `json:"tenantID"`
		BranchID     int64     `json:"branchID"`
		ProductID    int64     `json:"productID"`
		QuantitySold int       `json:"quantitySold"`
		PricePerUnit string    `json:"pricePerUnit"`
		Timestamp    time.Time `json:"timestamp"`
	}

	err := c.BindJSON(&input)
	if err != nil {
		app.errors.badRequestResponse(c, readJSONError(err))
		return
	}

	ppu, err := decimal.NewFromString(input.PricePerUnit)
	if err != nil {
		app.errors.badRequestResponse(c, err)
		return
	}

	txn := &data.SalesTransaction{
		TenantID:     input.TenantID,
		BranchID:     input.BranchID,
		ProductID:    input.ProductID,
		QuantitySold: input.QuantitySold,
		PricePerUnit: ppu,
		Timestamp:    input.Timestamp,
	}

	// validate
	v := validator.New()

	if data.ValidateSalesTransaction(v, txn); !v.Valid() {
		app.errors.failedValidationResponse(c, v.Errors)
		return
	}

	// aggregate transaction
	txnsChan <- *txn

	// response
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "Transaction received for processing",
	})
}

func (app *application) getSalesPerProduct(c *gin.Context) {
	// parse
	tenantID, err := strconv.ParseInt(c.Param("tenantID"), 10, 64)
	if err != nil {
		app.errors.badRequestResponse(c, err)
		return
	}

	// validate
	v := validator.New()
	v.Check(tenantID > 0, "tenantID", "invalid tenantID")

	if !v.Valid() {
		app.errors.failedValidationResponse(c, v.Errors)
		return
	}

	// query
	salesPerProduct, err := app.models.SalesTransactions.GetSalesPerProduct(tenantID)
	if err != nil {
		app.errors.serverErrorResponse(c, err)
		return
	}

	// response
	c.IndentedJSON(http.StatusCreated, gin.H{
		"sales_per_product": salesPerProduct,
	})
}

func (app *application) getTopSellingProducts(c *gin.Context) {
	// parse
	l := c.Query("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		app.errors.badRequestResponse(c, fmt.Errorf("unable to parse %v", l))
		return
	}

	// validate
	v := validator.New()
	v.Check(limit >= 1, "limit", "limit must be at least 1")
	v.Check(limit <= 100, "limit", "limit must be at most 100")

	if !v.Valid() {
		app.errors.failedValidationResponse(c, v.Errors)
		return
	}

	// query
	salesPerProduct, err := app.models.SalesTransactions.GetTopSellingProducts(limit)
	if err != nil {
		app.errors.serverErrorResponse(c, err)
		return
	}

	// response
	c.IndentedJSON(http.StatusCreated, gin.H{
		"sales_per_product": salesPerProduct,
	})
}
