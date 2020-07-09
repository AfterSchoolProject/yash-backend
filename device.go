package main

import "time"

type Device struct {
  ID          uint `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
  UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt   *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Name        string `json:"name" validate:"required"`
  Description string `json:"description,omitempty"`
  Host        string `json:"host" validate:"required,hostname"`
  Port        string `json:"port" validate:"required"`
  Actions     []Action `json:"actions"`
}
