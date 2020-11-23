package main

import (
  "time"

  "github.com/go-playground/validator/v10"
)

type User struct {
  ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at"`
  UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt          *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Login              string     `json:"login" validate:"required"`
  EncryptedPassword  string     `json:"password" validate:"required"`
  Devices            []Device   `json:"devices"`
}

func (user *User) Create() (err error) {

  if err = validator.New().Struct(user); err != nil {
    return
  }

  err = DB.Create(&user).Error
  return
}

func (user *User) CreateDevice(device *Device) (err error) {
  if err = validator.New().Struct(device); err != nil {
    return
  }

  err = DB.Model(&user).Association("Devices").Append(device).Error
  return
}
