package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var redisInstance *redis.Client

func InitRedis() (err error) {

	user := os.Getenv("REDIS_USERNAME")
	pwd := os.Getenv("REDIS_PASSWORD")
	addr := os.Getenv("REDIS_ADDRESS")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})

	val, err := rdb.Get(ctx, "count").Result()
	if err != nil {
		err = rdb.Set(ctx, "count", "0", 0).Err()
		if err != nil {
			//panic(err)
		}
	}
	fmt.Println("count = ", val)
	redisInstance = rdb
	return err
}

func GetRedis() *redis.Client {
	return redisInstance
}
