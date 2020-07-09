package main

import (
  "fmt"
  "encoding/json"
  "log"
  "net/http"
  "net/url"

  "github.com/jinzhu/gorm"
  "github.com/gorilla/mux"
  "github.com/go-playground/validator/v10"
)

func DevicesHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    log.Println("GET REQUEST")
    var devices []Device
    // var actions []Action
    DB.Preload("Actions").Find(&devices)
    for i := 0; i < len(devices); i++ {
      device := devices[i]
      if device.Actions == nil {
        device.Actions = []Action{}
      }
    }

    json.NewEncoder(w).Encode(devices)

  case http.MethodPost:
    log.Println("POST REQUEST")
    var device Device
    defer r.Body.Close()
    json.NewDecoder(r.Body).Decode(&device)

    validate := validator.New()
    err := validate.Struct(&device)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if err := DB.Create(&device).Error; err != nil {
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
  vars := mux.Vars(r)
  var device Device
  DB.Preload("Actions").First(&device, vars["id"])

  switch r.Method {
  case http.MethodGet:
    log.Println("GET DEVICE")
    json.NewEncoder(w).Encode(&device)

  case http.MethodPut:
    log.Println("UPDATE DEVICE")
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
    log.Println("DELETE DEVICE")

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

func postDevicesSendMessage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  var device Device
  DB.First(&device, vars["id"])

  // build request url
  post_url := url.URL{
    Scheme: "http",
    Host: device.Host + ":" + device.Port,
    Path: "action",
  }

  defer r.Body.Close()
  _, err := http.Post(post_url.String(), "application/json", r.Body)
  if err != nil {
    fmt.Println(err)
    http.Error(w, err.Error(), 500)
    return
  }

  // return response 200 OK
  w.WriteHeader(200)
}

func ActionsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  switch r.Method {
  case http.MethodGet:
    var device Device
    DB.Preload("Actions").First(&device, vars["id"])

    json.NewEncoder(w).Encode(device.Actions)

  case http.MethodPost:
    log.Println("ADD ACTION")
    var resp Action

    log.Println("DECODING BODY")
    json.NewDecoder(r.Body).Decode(&resp)
    log.Println(resp)

    var device Device
    DB.First(&device, vars["id"])
    if err := DB.Model(&device).Association("Actions").Append(&resp).Error; err != nil {
      log.Println("ERROR OCCURRED")
      log.Println(err)
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(201)
    log.Println(resp)
    json.NewEncoder(w).Encode(resp)

  default:
    w.WriteHeader(404)
  }
}


