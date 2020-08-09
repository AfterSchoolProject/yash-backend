package main

import "time"

type Action struct {
  ID          uint       `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
  UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt   *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Name        string     `json:"name" validate:"required"`
  Description string     `json:"description,omitempty"`
  Value       string     `json:"value"`
  DeviceID    uint       `json:"-"`
}
