package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahmad-abuziad/transaction-processor/internal/data"
	"github.com/ahmad-abuziad/transaction-processor/internal/metrics"
	"github.com/redis/go-redis/v9"
)

type SalesCache struct {
	Client  *redis.Client
	Metrics metrics.Metrics
}

func (c *SalesCache) SetSalesPerProduct(tenantID int64, salesPerProduct []data.SalesPerProduct) error {
	json, err := json.Marshal(salesPerProduct)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("sales_per_product:%d", tenantID)
	return c.Client.Set(context.Background(), key, string(json), 1*time.Minute).Err()
}

func (c *SalesCache) GetSalesPerProduct(tenantID int64) ([]data.SalesPerProduct, error) {
	cacheKey := fmt.Sprintf("sales_per_product:%d", tenantID)

	cachedValue, err := c.Client.Get(context.Background(), cacheKey).Result()
	if err != nil {
		c.Metrics.CacheMisses.Inc()
		return nil, err
	}
	c.Metrics.CacheHits.Inc()

	var salesPerProduct []data.SalesPerProduct
	err = json.Unmarshal([]byte(cachedValue), &salesPerProduct)
	if err != nil {
		return nil, err
	}

	return salesPerProduct, nil
}

const (
	TopSellingProductsKey        = "top_selling_products"
	TopSellingProductsExpiration = 10 * time.Minute
)

func (c *SalesCache) SetTopSellingProducts(topSellingProducts []data.TopSellingProduct) error {
	topProductsJSON, err := json.Marshal(topSellingProducts)
	if err != nil {
		return err
	}

	return c.Client.Set(context.Background(), TopSellingProductsKey, string(topProductsJSON), TopSellingProductsExpiration).Err()
}

func (c *SalesCache) GetTopSellingProducts() ([]data.TopSellingProduct, error) {
	cachedValue, err := c.Client.Get(context.Background(), TopSellingProductsKey).Result()
	if err != nil {
		c.Metrics.CacheMisses.Inc()
		return nil, err
	}
	c.Metrics.CacheHits.Inc()

	var topSellingProducts []data.TopSellingProduct
	err = json.Unmarshal([]byte(cachedValue), &topSellingProducts)
	if err != nil {
		return nil, err
	}

	return topSellingProducts, nil
}
