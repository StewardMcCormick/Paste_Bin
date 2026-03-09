package appcache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	"github.com/redis/go-redis/v9"
)

type apiKeyCache struct {
	client *redis.Client
	wg     *sync.WaitGroup
	quite  chan struct{}
}

func NewAPIKeyCache(client *redis.Client) *apiKeyCache {
	return &apiKeyCache{client: client, wg: &sync.WaitGroup{}}
}

func (c *apiKeyCache) Set(ctx context.Context, key string, value *domain.APIKey) {
	log := appctx.GetLogger(ctx)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		jsonValue, err := json.Marshal(value)
		if err != nil {
			log.Error(fmt.Sprintf("JSON parsing error - %v", err))
		}

		if c.client.Set(ctx, key, jsonValue, value.ExpiresAt.Sub(time.Now())).Err() != nil {
			log.Error(fmt.Sprintf("Redis saving error - %v", err))
		}
	}()
}

func (c *apiKeyCache) Get(ctx context.Context, key string) *domain.APIKey {
	log := appctx.GetLogger(ctx)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		log.Debug("cache miss")
		return nil
	}

	var result domain.APIKey
	if err := json.Unmarshal(data, &result); err != nil {
		log.Error(fmt.Sprintf("JSON parsing error - %v", err))
		return nil
	}

	log.Debug("cache hit")
	return &result
}

func (c *apiKeyCache) DeleteByKey(ctx context.Context, key string) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		log := appctx.GetLogger(ctx)
		if err := c.client.Del(ctx, key).Err(); err != nil {
			log.Error(fmt.Sprintf("Redis delete error - %v", err))
		}
	}()
}

func (c *apiKeyCache) Close(ctx context.Context) {
	log := appctx.GetLogger(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	close(c.quite)

	done := make(chan struct{})

	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		log.Error(fmt.Sprintf("Cache closing error - %v", ctx.Err()))
		return
	}
}
