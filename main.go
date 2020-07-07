package main

import (
  "fmt"
  "encoding/json"
  "log"
  "net/http"
  "net/url"
  "time"

  "github.com/rs/cors"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/go-playground/validator/v10"
)

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

type Action struct {
  ID          uint `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;not null"`
  CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
  UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
  DeletedAt   *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
  Name        string `json:"name" validate:"required"`
  Description string `json:"description,omitempty"`
  Value       string `json:"value"`
  DeviceID    uint `json:"-"`
}

var db *gorm.DB

func DevicesHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    log.Println("GET REQUEST")
    var devices []Device
    // var actions []Action
    db.Preload("Actions").Find(&devices)
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

    if err := db.Create(&device).Error; err != nil {
      http.Error(w, err.Error(), 400)
      return
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
  db.Preload("Actions").First(&device, vars["id"])

  switch r.Method {
  case http.MethodGet:
    log.Println("GET DEVICE")
    json.NewEncoder(w).Encode(&device)

  case http.MethodPut:
    log.Println("UPDATE DEVICE")
    var resp Device

    json.NewDecoder(r.Body).Decode(&resp)

    if err := db.Model(&device).Updates(map[string]interface{}{
      "name": resp.Name,
      "description": resp.Description,
      "host": resp.Host,
      "port": resp.Port,
    }).Error; err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(200)

  case http.MethodDelete:
    log.Println("DELETE DEVICE")

    if err := db.Unscoped().Delete(&device).Error; err != nil {
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
  db.First(&device, vars["id"])

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
    db.Preload("Actions").First(&device, vars["id"])

    json.NewEncoder(w).Encode(device.Actions)

  case http.MethodPost:
    log.Println("ADD ACTION")
    var resp Action

    log.Println("DECODING BODY")
    json.NewDecoder(r.Body).Decode(&resp)
    log.Println(resp)

    var device Device
    db.First(&device, vars["id"])
    if err := db.Model(&device).Association("Actions").Append(resp).Error; err != nil {
      log.Println("ERROR OCCURRED")
      log.Println(err)
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(201)

  default:
    w.WriteHeader(404)
  }
}

func main() {
  var err error
  db, err = gorm.Open("mysql", "root:eWkqxScKnE@/yash_dev?charset=utf8mb4&parseTime=True&loc=Local")
  if err != nil {
    log.Panic(err)
  }
  defer db.Close()
  log.Println("Database Connection Established")

  db.LogMode(true)
  db.AutoMigrate(&Device{}, &Action{})
  log.Println("Finished Migration")

  r := mux.NewRouter()
  r.HandleFunc("/devices", DevicesHandler).Methods(http.MethodGet, http.MethodPost)
  r.HandleFunc("/devices/{id}", DeviceHandler).Methods(http.MethodPut, http.MethodGet, http.MethodDelete)
  // r.HandleFunc("/devices/{id}/send_message", postDevicesSendMessage).Methods("POST")

  r.HandleFunc("/devices/{id}/actions", ActionsHandler)

  handler := cors.Default().Handler(r)

  c := cors.New(cors.Options{
    AllowedMethods: []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
  })

  handler = c.Handler(handler)
  log.Println("Server listening on localhost port 8080")
  log.Fatal(http.ListenAndServe(":8080", handler))
}
