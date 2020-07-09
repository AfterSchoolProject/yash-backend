package main

import (
  "log"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func SetupDatabase() *gorm.DB {
  var err error
  DB, err = gorm.Open("mysql", "root:eWkqxScKnE@/yash_dev?charset=utf8mb4&parseTime=True&loc=Local")
  if err != nil {
    log.Panic(err)
  }
  log.Println("Database Connection Established")

  DB.LogMode(true)
  DB.AutoMigrate(&Device{}, &Action{})
  log.Println("Finished Migration")

  return DB
}
