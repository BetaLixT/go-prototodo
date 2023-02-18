package mcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewMemoryCache() *cache.Cache {
  return cache.New(
    time.Minute * 5,
    time.Minute * 10,
  )
}
