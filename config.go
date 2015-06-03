/*
    conntrack-logger
    Copyright (C) 2015 Denis V Chapligin <akashihi@gmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
    "flag"
    "os"
    "fmt"
    "encoding/json"
)

type Configuration struct {
    Workers     int
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
