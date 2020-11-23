package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "os"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "gopkg.in/yaml.v3"
)

var DB *gorm.DB

type DBConfig struct {
  User     string
  Password string
  Database string
}

func SetupDatabase() *gorm.DB {
  data, err := ioutil.ReadFile("database.yml")
  if err != nil {
    fmt.Printf("Error reading database configs %s", err)
    return nil
  }

  env := os.Getenv("GO_ENV")
  if env == "" {
    env = "develop"
  }

  yml := map[string]DBConfig{}
  yaml.Unmarshal(data, &yml)

  config   := yml[env]
  user     := config.User
  password := config.Password
  database := config.Database
  dbUrl    := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)

  DB, err = gorm.Open("mysql", dbUrl)
  if err != nil {
    log.Panic(err)
  }
  log.Println("Database Connection Established")

  DB.LogMode(false)
  DB.AutoMigrate(&Device{}, &Action{}, &User{})
  log.Println("Finished Migration")

  return DB
}
