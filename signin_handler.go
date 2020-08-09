package main

import (
  "encoding/json"
  "net/http"

  "golang.org/x/crypto/bcrypt"
)

type Credentials struct {
  Login    string `json:"login"`
  Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
  var credentials Credentials

  json.NewDecoder(r.Body).Decode(&credentials)

  var user User
  query := User{Login: credentials.Login}

  if DB.Where(&query).First(&user).RecordNotFound() {
    http.Error(w, "User Not Found", 403)
    return
  }

  err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(credentials.Password))
  if err != nil {
    http.Error(w, "UnAuthorized", 403)
    return
  }

  session, _ := store.Get(r, "session")
  session.Values["login"] = user.Login
  session.Values["Authenticated"] = true

  if err = session.Save(r, w); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)

    return
  }
}
