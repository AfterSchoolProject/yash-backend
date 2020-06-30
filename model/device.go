package model

import (
  "log"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
  database, err := gorm.Open("mysql", "root:eWkqxScKnE@/yash_dev?charset=utf8mb4&parseTime=True&loc=Local")

  if err != nil {
    log.Panic(err)
  }

  database.AutoMigrate(&Device{})

  DB = database
}
