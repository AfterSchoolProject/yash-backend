package main

import "time"

type User struct {
  ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at"`
  UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt          *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Login              string     `json:"login" validate:"required,unique"`
  EncryptedPassword  string     `json:"password" validate:"required"`
  Devices            []Device   `json:"devices"`
}
