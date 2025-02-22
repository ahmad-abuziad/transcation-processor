package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	gin := gin.Default()
	gin.GET("/health", health)
	gin.POST("/tenant/:tenantID/branch/:branchID/sales-transaction", app.newSalesTransaction)
	gin.GET("/tenant/:tenantID/sales", app.getSalesPerProduct)
	gin.GET("/sales", app.getTopSellingProducts)

	return gin
}
