package main

import (
  "time"

  "github.com/go-playground/validator/v10"
)

type Device struct {
  ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
  UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt   *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Name        string     `json:"name" validate:"required"`
  Description string     `json:"description,omitempty"`
  Host        string     `json:"host" validate:"required,hostname_rfc1123"`
  Port        string     `json:"port" validate:"required"`
  UserID      uint       `json:"-"`
  Actions     []Action   `json:"actions"`
}

func (device *Device) CreateAction(action *Action) (err error) {
  if err = validator.New().Struct(action); err != nil {
    return
  }

  err = DB.Model(&device).Association("Actions").Append(&action).Error
  return
}

func (device *Device) Delete() (err error) {
  err = DB.Unscoped().Select("Actions").Delete(&device).Error

  return
}
