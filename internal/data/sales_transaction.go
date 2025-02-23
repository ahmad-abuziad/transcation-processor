package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ahmad-abuziad/transaction-processor/internal/validator"
	"github.com/shopspring/decimal"
)

type SalesTransaction struct {
	ID           int64
	TenantID     int64
	BranchID     int64
	ProductID    int64
	QuantitySold int
	PricePerUnit decimal.Decimal
	Timestamp    time.Time
}

func ValidateSalesTransaction(v *validator.Validator, txn *SalesTransaction) {
	v.Check(txn.TenantID > 0, "tenantID", "invalid identifier")
	v.Check(txn.BranchID > 0, "branchID", "invalid identifier")
	v.Check(txn.ProductID > 0, "productID", "invalid identifier")
	v.Check(txn.QuantitySold > 0, "quantitySold", "must be greater than 0")
	v.Check(txn.PricePerUnit.IsPositive(), "pricePerUnit", "must be greater than 0")
	v.Check(txn.Timestamp.Before(time.Now().Add(1*time.Second)), "timestamp", "must not be in the future")
}

type SalesTransactionModel struct {
	DB *sql.DB
}

func (m *SalesTransactionModel) Insert(txn *SalesTransaction) error {
	stmt := `INSERT INTO sales_transactions (tenant_id, branch_id, product_id, quantity_sold, price_per_unit, log_timestamp)
    VALUES(?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, txn.TenantID, txn.BranchID, txn.ProductID, txn.QuantitySold, txn.PricePerUnit, txn.Timestamp)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	txn.ID = id

	return nil
}

type SalesPerProduct struct {
	ProductID  int64
	TotalSales decimal.Decimal
}

func (m *SalesTransactionModel) GetSalesPerProduct(tenantID int64) ([]SalesPerProduct, error) {
	stmt := `SELECT product_id, SUM(quantity_sold * price_per_unit) AS total_sales
	FROM sales_transactions
	WHERE tenant_id = ?
	GROUP BY product_id
	ORDER BY total_sales DESC;`

	rows, err := m.DB.Query(stmt, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	spp := []SalesPerProduct{}

	for rows.Next() {
		var s SalesPerProduct

		err = rows.Scan(&s.ProductID, &s.TotalSales)
		if err != nil {
			return nil, err
		}

		spp = append(spp, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return spp, nil
}

type TopSellingProduct struct {
	ProductID     int64
	TotalQuantity int
}

func (m *SalesTransactionModel) GetTopSellingProducts(limit int) ([]TopSellingProduct, error) {
	stmt := `SELECT product_id, SUM(quantity_sold) AS total_quantity FROM sales_transactions
	GROUP BY product_id
	ORDER BY total_quantity DESC
	LIMIT ?;`

	rows, err := m.DB.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tsp := []TopSellingProduct{}

	for rows.Next() {
		var s TopSellingProduct

		err = rows.Scan(&s.ProductID, &s.TotalQuantity)
		if err != nil {
			return nil, err
		}

		tsp = append(tsp, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tsp, nil
}

func (m *SalesTransactionModel) InsertBatch(txns []SalesTransaction) {
	if len(txns) == 0 {
		return
	}

	query := "INSERT INTO sales_transactions (tenant_id, branch_id, product_id, quantity_sold, price_per_unit, log_timestamp) VALUES "
	values := []interface{}{}

	for _, txn := range txns {
		query += "(?, ?, ?, ?, ?, ?),"
		values = append(values, txn.TenantID, txn.BranchID, txn.ProductID, txn.QuantitySold, txn.PricePerUnit, txn.Timestamp)
	}

	query = query[:len(query)-1] // Remove last comma

	_, err := m.DB.Exec(query, values...)
	if err != nil {
		// TODO handle batch failure, retry, log failed transactions after retries
		fmt.Printf("Failed to insert batch: %v\n", err)
	}
}
