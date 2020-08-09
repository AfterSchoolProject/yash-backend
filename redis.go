package main

import (
  "context"
  "log"

  "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var reds *redis.Client

func SetupRedis() {
  reds = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
  })

  pong, err := reds.Ping(ctx).Result()
  if err != nil {
    log.Panic(pong, err)
  }

  log.Println("Redis server started")
}
