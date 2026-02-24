package conn

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(redisServer string) (*redis.Client) {

    // Open a connection to the redis
	RedisClient := redis.NewClient(&redis.Options{
        Addr: redisServer,
	})

	// Verify the connection is established and valid by pinging redis
    _, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal("❌ Failed to connect to Redis:", err)
    }

	log.Println("✅ Successfully connected to Redis!")

	return RedisClient
}