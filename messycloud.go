package main

import (
  "stathat.com/c/jconfig"
  "os"
  "log"
)

func main() {
/*
  // Prod settings for logging
  if f, err := os.OpenFile("/var/log/messycloud.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644); err != nil {
    log.Fatalf("error opening file: %v\n", err)
  } else {
    defer f.Close()
    log.SetOutput(f)
  }
*/

  config := jconfig.LoadConfig("./settings.json")
  datafolder := config.GetString("datafolder")
  if datafolder == "" {
    log.Fatalln("datafolder not specifed or empty in settings.json")
  }

  watchFiles(datafolder)
}