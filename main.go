package main

func main() {
  db := SetupDatabase()
  defer db.Close()

  SetupRoutesAndRun()
}
