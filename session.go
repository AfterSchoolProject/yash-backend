package main

import (
  "log"
  "os"

  redisStore "gopkg.in/boj/redistore.v1"
  "github.com/gorilla/sessions"
)

var store *redisStore.RediStore

func init() {
  var err error
  store, err = redisStore.NewRediStore(10, "tcp", ":6379", "", []byte(os.Getenv("SESSION_KEY")))
  if err != nil {
    log.Panic(err)
  }

  store.Options = &sessions.Options{
    MaxAge: 86400, // 1 day
    HttpOnly: true,
  }

  log.Println("Session store initialized")
}
