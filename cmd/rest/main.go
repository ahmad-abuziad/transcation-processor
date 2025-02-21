package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var count int

func main() {
	r := gin.Default()
	r.GET("/health", health)

	// use Zap or Logrus
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(os.Getenv("DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS message (value INT AUTO_INCREMENT PRIMARY KEY)"); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	_, err = db.Exec("INSERT INTO message VALUES ()")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	c, err := countRecords(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	count = c

	r.Run()
}

func countRecords(db *sql.DB) (int, error) {

	rows, err := db.Query("SELECT COUNT(*) FROM message")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		rows.Close()
	}

	return count, nil
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Count":  count,
	})
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
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
