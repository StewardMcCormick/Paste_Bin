package redis

import "time"

type RedisConfig interface {
	GetHost() string
	GetPort() int
	GetPoolSize() int
	GetPoolTimeout() time.Duration
	GetPassword() string
	GetDb() int
}

func (c CacheConfig) GetHost() string {
	return c.Host
}

func (c CacheConfig) GetPort() int {
	return c.Port
}

func (c CacheConfig) GetPoolSize() int {
	return c.PoolSize
}

func (c CacheConfig) GetPoolTimeout() time.Duration {
	return c.PoolTimeout
}

func (c CacheConfig) GetPassword() string {
	return c.Password
}

func (c CacheConfig) GetDb() int {
	return c.Db
}

func (c RateConfig) GetHost() string {
	return c.Host
}

func (c RateConfig) GetPort() int {
	return c.Port
}

func (c RateConfig) GetPoolSize() int {
	return c.PoolSize
}

func (c RateConfig) GetPoolTimeout() time.Duration {
	return c.PoolTimeout
}

func (c RateConfig) GetPassword() string {
	return c.Password
}

func (c RateConfig) GetDb() int {
	return c.Db
}
