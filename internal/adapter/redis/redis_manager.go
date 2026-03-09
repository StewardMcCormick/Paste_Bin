package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	apiKeyCacheName = "API_KEY_CACHE"
	pasteCacheName  = "PASTE_CACHE"
	rateCacheName   = "RATE_LIMITING"
)

type CacheConfig struct {
	Host        string        `env:"REDIS_CACHE_HOST" required:"true"`
	Port        int           `env:"REDIS_CACHE_PORT" required:"true"`
	PoolSize    int           `yaml:"pool_size" env-default:"10"`
	PoolTimeout time.Duration `yaml:"pool_timeout" env-default:"5s"`
	User        string        `env:"REDIS_CACHE_USER" required:"true"`
	Password    string        `env:"REDIS_CACHE_PASSWORD" required:"true"`
	Db          int
}

type RateConfig struct {
	Host        string        `env:"REDIS_RATE_HOST" required:"true"`
	Port        int           `env:"REDIS_RATE_PORT" required:"true"`
	PoolSize    int           `yaml:"pool_size" env-default:"10"`
	PoolTimeout time.Duration `yaml:"pool_timeout" env-default:"5s"`
	User        string        `env:"REDIS_RATE_USER" required:"true"`
	Password    string        `env:"REDIS_RATE_PASSWORD" required:"true"`
	Db          int
}

type Config struct {
	Cache *CacheConfig
	Rate  *RateConfig
}

type Client struct {
	*redis.Client
	name string
}

type Manager struct {
	clients map[string]*Client
	mu      *sync.Mutex
}

func NewManager(cfg Config) (*Manager, error) {
	manager := &Manager{
		clients: make(map[string]*Client, 2),
		mu:      &sync.Mutex{},
	}

	if cfg.Cache != nil {
		cfg.Cache.Db = 0
		apiKeyCacheClient, err := newClient(apiKeyCacheName, cfg.Cache)
		if err != nil {
			return nil, err
		}

		cfg.Cache.Db = 1
		pasteCacheClient, err := newClient(pasteCacheName, cfg.Cache)
		if err != nil {
			return nil, err
		}

		manager.clients[apiKeyCacheClient.name] = apiKeyCacheClient
		manager.clients[pasteCacheClient.name] = pasteCacheClient
	}

	if cfg.Rate != nil {
		cfg.Rate.Db = 0
		rateClient, err := newClient(rateCacheName, cfg.Rate)
		if err != nil {
			return nil, err
		}

		manager.clients[rateClient.name] = rateClient
	}

	return manager, nil
}

func newClient[T RedisConfig](name string, cfg T) (*Client, error) {

	addr := fmt.Sprintf("%s:%d", cfg.GetHost(), cfg.GetPort())

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Username:    cfg.GetUser(),
		Password:    cfg.GetPassword(),
		PoolSize:    cfg.GetPoolSize(),
		PoolTimeout: cfg.GetPoolTimeout(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis client init error: %s, error: %w", name, err)
	}

	return &Client{
		client,
		name,
	}, nil
}

func (m *Manager) GetAPIKeyCacheClient() *Client {
	return m.clients[apiKeyCacheName]
}

func (m *Manager) GetPasteCacheClient() *Client {
	return m.clients[pasteCacheName]
}

func (m *Manager) Close() error {
	var err error
	for _, client := range m.clients {
		err = client.Close()
	}

	return err
}
