package cache

import (
	"github.com/ahmad-abuziad/transaction-processor/internal/metrics"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	SalesCache SalesCache
}

func NewCache(client *redis.Client, metrics metrics.Metrics) Cache {
	return Cache{
		SalesCache: SalesCache{Client: client, Metrics: metrics},
	}
}
