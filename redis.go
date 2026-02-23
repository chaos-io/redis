package redis

import (
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	redisClient     goredis.UniversalClient
	redisClientOnce sync.Once
)

type Redis struct {
	goredis.UniversalClient
}

func GetRedis() goredis.UniversalClient {
	redisClientOnce.Do(func() {
		cfg := NewConfig()
		redisClient = New(cfg)
	})
	return redisClient
}

func New(cfg *Config) *Redis {
	if cfg == nil {
		cfg = NewConfig()
	}

	if len(cfg.Connections) == 0 {
		cfg.Connections = []string{":6379"}
	}

	if cfg.MinIdleConns == 0 {
		if cfg.MaxIdleConns > 0 {
			cfg.MinIdleConns = cfg.MaxIdleConns
		} else {
			cfg.MinIdleConns = 100
		}
	}

	if cfg.PoolSize == 0 {
		cfg.PoolSize = 300
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = cfg.ReadTimeout
	}

	if len(cfg.Connections) == 1 {
		option := &goredis.Options{
			Addr:            cfg.Connections[0],
			Password:        cfg.Password,
			DB:              cfg.DB,
			MinIdleConns:    cfg.MinIdleConns,
			PoolSize:        cfg.PoolSize,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			MaxRetries:      cfg.MaxRetries,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			MinRetryBackoff: cfg.MinRetryBackoff,
		}
		return &Redis{goredis.NewClient(option)}
	}

	clusterOption := &goredis.ClusterOptions{
		Addrs:           cfg.Connections,
		Password:        cfg.Password,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxRetries:      cfg.MaxRetries,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		MinRetryBackoff: cfg.MinRetryBackoff,
	}
	return &Redis{goredis.NewClusterClient(clusterOption)}
}
