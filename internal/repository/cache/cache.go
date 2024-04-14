package cache

import (
	"avito-banner/internal/domain"
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Cache struct {
	Banners map[Ids]domain.CacheUserBanner
	mx      sync.RWMutex
}

type Ids struct {
	TagId     int
	FeatureId int
}

func NewCache() *Cache {
	return &Cache{Banners: make(map[Ids]domain.CacheUserBanner)}
}

func (c *Cache) Get(ctx context.Context, tagId, featureId int) (domain.CacheUserBanner, error) {
	log.Println("Start Get from cache")
	defer log.Println("End Get from cache")

	ids := Ids{tagId, featureId}
	c.mx.RLock()

	value, exist := c.Banners[ids]
	c.mx.RUnlock()
	if !exist {
		return domain.CacheUserBanner{}, errors.New("key doesn't exist in cache")
	}

	if time.Now().After(value.Expiration) {
		deleteFromCache(c, ids)
		return domain.CacheUserBanner{}, errors.New("key expired")
	}

	return value, nil
}

func (c *Cache) Set(ctx context.Context, tagId, featureId int, title, text, url string, isActive bool) error {
	log.Println("Start Set for cache")
	defer log.Println("End Set for cache")

	ids := Ids{tagId, featureId}
	data := domain.CacheUserBanner{
		Title:      title,
		Text:       text,
		Url:        url,
		IsActive:   isActive,
		Expiration: time.Now().Add(5 * time.Minute),
	}
	c.mx.Lock()
	c.Banners[ids] = data
	c.mx.Unlock()
	return nil
}

func deleteFromCache(c *Cache, key Ids) {
	c.mx.Lock()
	delete(c.Banners, key)
	c.mx.Unlock()
}

func (c *Cache) BackgroundCleaning(ctx context.Context) {
	defer log.Println("Complete cache background cleaning")
	for {
		time.Sleep(1 * time.Minute)
		select {
		case <-ctx.Done():
			return

		default:
			now := time.Now()
			c.mx.RLock()
			for key, value := range c.Banners {
				c.mx.RUnlock()
				if now.After(value.Expiration) {
					deleteFromCache(c, key)
				}
				c.mx.RLock()
			}
			c.mx.RUnlock()
		}
	}
}
