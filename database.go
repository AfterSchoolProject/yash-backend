package main

import (
  "fmt"
  "log"
  "os"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func SetupDatabase() *gorm.DB {
  user     := os.Getenv("DATABASE_USER")
  password := os.Getenv("DATABASE_PASSWORD")
  database := os.Getenv("DATABASE")
  dbUrl    := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)

  var err error
  DB, err = gorm.Open("mysql", dbUrl)
  if err != nil {
    log.Panic(err)
  }
  log.Println("Database Connection Established")

  DB.LogMode(true)
  DB.AutoMigrate(&Device{}, &Action{}, &User{})
  log.Println("Finished Migration")

  return DB
}
