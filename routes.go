package main

import (
  "log"
  "net/http"

  "github.com/rs/cors"
  "github.com/gorilla/mux"
)

func SetupRoutesAndRun() {
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
