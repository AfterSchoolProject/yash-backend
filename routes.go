package main

import (
  "log"
  "net/http"

  "github.com/rs/cors"
  "github.com/gorilla/mux"
)

func Authenticate(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session")
    if session.Values["Authenticated"] != true {
      log.Printf("User %s not authenticated!", session.Values["login"])
      http.Error(w, "Not Authenticated", 403)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func SetupRoutesAndRun() {
  r := mux.NewRouter()

  r.HandleFunc("/signin", SignInHandler).Methods(http.MethodPost)

  s := r.PathPrefix("/devices").Subrouter()
  s.Use(Authenticate)

  // Devices 
  s.HandleFunc("", DevicesHandler).Methods(http.MethodGet, http.MethodPost)
  s.HandleFunc("/{id}", DeviceHandler).Methods(http.MethodPut, http.MethodGet, http.MethodDelete)

  // Actions
  s.HandleFunc("/{id}/actions", ActionsHandler).Methods(http.MethodGet, http.MethodPost)
  s.HandleFunc("/{device_id}/actions/{id}", ActionHandler).Methods(http.MethodPost, http.MethodDelete)

  c := cors.New(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000"},
    AllowedMethods: []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
    AllowCredentials: true,
    Debug: true,
  })

  handler := c.Handler(r)

  log.Println("Server listening on localhost port 8080")
  log.Fatal(http.ListenAndServe(":8080", handler))
}
