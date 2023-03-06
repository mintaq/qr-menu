package cache

import (
	"os"
	"strconv"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// RedisConnection func for connect to Redis server.
func RedisConnection() {
	// Define Redis database number.
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("REDIS_MAX_IDLE_CONNECTIONS"))

	// Build Redis connection URL.
	redisConnURL, err := utils.ConnectionURLBuilder("redis")
	if err != nil {
		panic(err)
	}

	// Set Redis options.
	options := &redis.Options{
		MaxIdleConns: maxIdleConns,
		Addr:         redisConnURL,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           dbNumber,
	}

	rdb := redis.NewClient(options)

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	RedisClient = rdb
}
