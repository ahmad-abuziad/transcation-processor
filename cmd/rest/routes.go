package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() http.Handler {
	r := gin.Default()

	r.Use(app.requestDurationMiddleware())

	r.GET("/health", health)

	r.POST("/tenant/:tenantID/branch/:branchID/sales-transaction", app.newSalesTransaction)
	r.GET("/tenant/:tenantID/sales", app.getSalesPerProduct)
	r.GET("/top-selling", app.getTopSellingProducts)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func (app *application) requestDurationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		app.metrics.RequestDuration.WithLabelValues(c.FullPath()).Observe(duration)
	}
}
