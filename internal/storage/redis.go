package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func NewRedis() (*redis.Client, error) {
	options := &redis.Options{
		Addr:        os.Getenv("REDIS_HOST"),
		Password:    "",
		DB:          0,
		ReadTimeout: 2 * time.Second,
	}
	rds := redis.NewClient(options)

	fmt.Println(rds.Ping().Result())
	return rds, nil
}
