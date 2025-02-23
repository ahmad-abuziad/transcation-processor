package cache

import (
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	SalesCache SalesCache
}

func NewCache(client *redis.Client) Cache {
	return Cache{
		SalesCache: SalesCache{Client: client},
	}
}
