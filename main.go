package main

import (
  "fmt"
  "encoding/json"
  "log"
  "net/http"
  "net/url"
  "time"

  // "github.com/gorilla/handlers"
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
  Name        string `json:"name" gorm:"unique" validate:"required"`
  Description string `json:"description,omitempty"`
  Host        string `json:"host" validate:"required,hostname"`
  Port        string `json:"port" validate:"required"`
}

var db *gorm.DB

func handleDevices(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")

  switch r.Method {
  case http.MethodOptions:
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    w.WriteHeader(200)
  case http.MethodGet:
    log.Println("GET REQUEST")
    var devices []Device
    db.Find(&devices)

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

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(device)

  default:
    w.WriteHeader(404)
  }
  // w.Header().Set("Access-Control-Allow-Origin", "*")
}

func postDevicesSendMessage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  var device Device
  db.First(&device, vars["id"])

  w.Header().Set("Access-Control-Allow-Origin", "*")

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

func main() {
  var err error
  db, err = gorm.Open("mysql", "root:eWkqxScKnE@/yash_dev?charset=utf8mb4&parseTime=True&loc=Local")
  if err != nil {
    log.Panic(err)
  }
  defer db.Close()
  log.Println("Database Connection Established")

  db.AutoMigrate(&Device{})
  log.Println("Finished Migration")

  r := mux.NewRouter()
  r.HandleFunc("/devices", handleDevices)
  // r.HandleFunc("/devices/{id}/send_message", postDevicesSendMessage).Methods("POST")
  r.Use(mux.CORSMethodMiddleware(r))

  log.Println("Server listening on localhost port 8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}
