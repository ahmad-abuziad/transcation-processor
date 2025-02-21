package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
}

func main() {

	app := application{}

	// use Zap or Logrus
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(os.Getenv("DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app.serve()
}

func (app *application) serve() *gin.Engine {
	r := gin.Default()
	r.GET("/health", health)
	r.POST("/tenant/:tenantID/branch/:branchID/sales-transaction", newSalesTransaction)
	r.GET("/tenant/:tenantID/sales", getSalesPerProduct)
	r.GET("/sales", getTopSellingProducts)

	return r
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
