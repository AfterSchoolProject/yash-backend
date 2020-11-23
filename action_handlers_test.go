package main

import (
  "bytes"
  "errors"
  "fmt"
  "encoding/json"
  "net"
  "net/http"
  "net/http/httptest"
  "net/url"
  "os"
  "testing"

  "golang.org/x/crypto/bcrypt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
  SetupSession()
  os.Exit(m.Run())
}

func AuthorizedUser() (user *User) {
  SetupDatabase()

  user = &User{}
  result := DB.Where(&User{Login: "test"}).First(&user)

  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    encrypted, _ := bcrypt.GenerateFromPassword([]byte("nothing"), 10)
    user = &User{Login: "test", EncryptedPassword: string(encrypted)}
    if err := user.Create(); err != nil {
      fmt.Printf("User create failed %s\n", err)
      panic("USER CREATE FAILED")
    }
  }

  return
}

func ClearTables() {
  DB.Exec("DELETE FROM devices")
  DB.Exec("DELETE FROM actions")
}

func SetupSession() (err error) {
  user := AuthorizedUser()
  r, _ := http.NewRequest("GET", "localhost", nil)
  w := httptest.NewRecorder()
  fmt.Printf("User session: %v\n", user)
  session, _ := store.Get(r, "session")
  session.Values["login"] = user.Login
  session.Values["Authenticated"] = true

  err = session.Save(r, w)
  return
}

func TestActionsHandler(t *testing.T) {
  user := AuthorizedUser()
  t.Run("all actions are returned for a device", func(t *testing.T) {
    ClearTables()

    device := Device{
      Name: "all actions device",
      Host: "localhost",
      Port: "3000",
      Actions: []Action{{Name: "all actions"}},
    }
    if err := user.CreateDevice(&device); err != nil {
      t.Fatalf("Device could not be created %s", err)
    }

    url := fmt.Sprintf("/devices/%d/actions", device.ID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      t.Errorf("Request failed")
    }

    rr := httptest.NewRecorder()
    router := mux.NewRouter()
    router.HandleFunc("/devices/{id}/actions", ActionsHandler)
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
      t.Errorf("handler returned wrong status coe: got %v want %v", status, http.StatusOK)
    }

    actions := []Action{}
    json.NewDecoder(rr.Result().Body).Decode(&actions)
    defer rr.Result().Body.Close()

    if len(actions) != 1 {
      t.Fatalf("Did not return all of device's actions: %d", len(actions))
    }

    action := actions[0]
    if action.Name != "all actions" {
      t.Errorf("Name did not match: %s", action.Name)
    }
  })

  t.Run("an action is created for a device", func(t *testing.T) {
    ClearTables()

    device := Device{
      Name: "creating actions device",
      Host: "localhost",
      Port: "3000",
    }
    if err := user.CreateDevice(&device); err != nil {
      t.Fatalf("Device could not be created %s", err)
    }

    var b bytes.Buffer
    requestBody := Action{Name: "created action"}
    json.NewEncoder(&b).Encode(requestBody)

    url := fmt.Sprintf("/devices/%d/actions", device.ID)
    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
      t.Errorf("Request failed")
    }

    rr := httptest.NewRecorder()

    router := mux.NewRouter()
    router.HandleFunc("/devices/{id}/actions", ActionsHandler)
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
      t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    responseBody := Action{}
    json.NewDecoder(rr.Result().Body).Decode(&responseBody)
    defer rr.Result().Body.Close()

    if responseBody.Name != "created action" {
      t.Errorf("action was not created with the right values %s", responseBody.Name)
    }
  })
}

func TestActionHandler(t *testing.T) {
  user := AuthorizedUser()
  t.Run("sends a request to external device", func(t *testing.T) {
    ClearTables()

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      var requestBody map[string]string
      json.NewDecoder(r.Body).Decode(&requestBody)
      defer r.Body.Close()

      if requestBody["value"] != "action value" {
        t.Fatalf("the wrong value was sent %s", requestBody["value"])
      }

      w.WriteHeader(200)
    }))
    defer ts.Close()

    u, _ := url.Parse(ts.URL)
    host, port, _ := net.SplitHostPort(u.Host)

    device := &Device{
      Name: "test device",
      Host: host,
      Port: port,
      Actions: []Action{{Name: "test action", Value: "action value"}},
    }
    if err := user.CreateDevice(device); err != nil {
      t.Fatalf("Device could not be created %s", err)
    }

    requestUrl := fmt.Sprintf("/devices/%d/actions/%d", device.ID, device.Actions[0].ID)
    req, err := http.NewRequest("POST", requestUrl, nil)
    if err != nil {
      t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    router := mux.NewRouter()
    router.HandleFunc("/devices/{device_id}/actions/{id}", ActionHandler)
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
      t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
  })
}
