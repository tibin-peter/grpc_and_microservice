package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis()*redis.Client{
	addr:=os.Getenv("REDIS_ADDR")

	rdb:=redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_,err:=rdb.Ping(Ctx).Result()
	if err!=nil{
		panic("failed to connect redis")
	}

	fmt.Println("User service connected to Redis")

	return rdb
}