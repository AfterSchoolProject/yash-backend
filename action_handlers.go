package main

import (
  "bytes"
  "fmt"
  "encoding/json"
  "log"
  "net/http"
  "net/url"
  "strconv"

  "github.com/gorilla/mux"
)

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

type Message struct {
  Value string `json:"value"`
}

func ActionHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  var device Device
  DB.First(&device, vars["device_id"])

  // build request url
  postUrl := url.URL{
    Scheme: "http",
    Host: device.Host + ":" + device.Port,
    Path: "action",
  }

  id, _ := strconv.ParseUint(vars["id"], 10, 64)
  action := Action{ID: uint(id)}
  DB.Model(&device).Association("Actions").Find(&action)

  jsonValue, err := json.Marshal(Message{Value: action.Value})
  requestBody := bytes.NewReader(jsonValue)

  defer r.Body.Close()
  _, err = http.Post(postUrl.String(), "application/json", requestBody)
  if err != nil {
    fmt.Println(err)
    http.Error(w, err.Error(), 500)

    return
  }

  w.WriteHeader(200)
}
