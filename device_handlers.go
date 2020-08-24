package main

import (
  "encoding/json"
  "log"
  "net/http"

  "github.com/jinzhu/gorm"
  "github.com/gorilla/mux"
  "github.com/go-playground/validator/v10"
)

func DevicesHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "session")

  login, ok := session.Values["login"].(string)
  if !ok {
    log.Println("ERROR WHEN GETING SESSION")
    http.Error(w, "Login not saved", 500)
    return
  }

  var user User
  DB.Where(User{Login: login}).First(&user)

  switch r.Method {
  case http.MethodGet:
    DB.Preload("Devices").Find(&user)

    for i := 0; i < len(user.Devices); i++ {
      device := &user.Devices[i]
      DB.Preload("Actions").Find(device)
    }

    json.NewEncoder(w).Encode(user.Devices)

  case http.MethodPost:
    var device Device
    defer r.Body.Close()
    json.NewDecoder(r.Body).Decode(&device)

    validate := validator.New()
    err := validate.Struct(&device)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if err := DB.Model(&user).Association("Devices").Append(&device).Error; err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if device.Actions == nil {
      device.Actions = []Action{}
    }

    w.WriteHeader(201)
    json.NewEncoder(w).Encode(device)

  default:
    w.WriteHeader(404)
  }
}

func DeviceHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "session")
  login, ok := session.Values["login"].(string)
  if !ok {
    http.Error(w, "Login not saved", 500)
  }

  var user User
  DB.Where(User{Login: login}).First(&user)

  vars := mux.Vars(r)
  var device Device

  DB.Where("user_id = ? AND id = ?", user.ID, vars["id"]).Preload("Actions").First(&device)

  switch r.Method {
  case http.MethodGet:
    json.NewEncoder(w).Encode(&device)

  case http.MethodPut:
    var resp Device

    json.NewDecoder(r.Body).Decode(&resp)

    if err := DB.Model(&device).Updates(map[string]interface{}{
      "name": resp.Name,
      "description": resp.Description,
      "host": resp.Host,
      "port": resp.Port,
    }).Error; err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&device)

  case http.MethodDelete:
    err := DB.Transaction(func(tx *gorm.DB) error {
      if err := tx.Unscoped().Where("device_id = ?", device.ID).Delete(&Action{}).Error; err != nil {
        return err
      }

      if err := tx.Unscoped().Delete(&device).Error; err != nil {
        return err
      }

      return nil
    })

    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(200)

  default:
    w.WriteHeader(404)
  }
}
