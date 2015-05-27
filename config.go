package main

import (
    "flag"
    "os"
    "fmt"
    "encoding/json"
)

type Configuration struct {
    CommitCount int
    DBhost      string
    DBport      string
    DBuser      string
    DBpassword  string
    DBname      string
}


func config() Configuration {
  //Command line options
  configFilePtr := flag.String("config", "/etc/conntrack-logger.cfg", "Path to the config file")
  flag.Parse()

  //Configuration
  file, e := os.Open(*configFilePtr)
  if e != nil {
    fmt.Println("Unable to open config file:", *configFilePtr)
    os.Exit(1)
  }
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  e = decoder.Decode(&configuration)
  if e != nil {
    fmt.Println("error:", e)
    os.Exit(1)
  }
  return configuration
}
