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
    _ "github.com/lib/pq"
    "database/sql"
    "log"
    "time"
)

const (
    PURGE_PERIOD = 24
)

type FlowRecord struct {
    TS       time.Time
    Proto    string
    Src      string
    Dst      string
    Sport    int
    Dport    int
}

func writeDB(configuration Configuration, flows <-chan FlowRecord) {
    connect_string := "host="+configuration.DBhost+" port="+configuration.DBport+" user="+configuration.DBuser+" password="+configuration.DBpassword+" dbname="+configuration.DBname+" sslmode=disable"
    

    for {
	db, err := sql.Open("postgres", connect_string)
	defer db.Close()

	if err != nil {
    	    log.Print(err)
            time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
            continue
	}
	log.Print("Connection to database established: "+configuration.DBhost)

	for {
	    counter := 0

	    tx, err := db.Begin()
	    if err != nil {
		log.Print(err)
                time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
		break
	    }

	    insert_query, err := tx.Prepare("insert into events (ts, proto, src, dst, sport, dport) values ($1, $2, $3, $4, $5, $6)")

	    if err != nil {
		log.Print(err)
        	time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
        	break
    	    }

	    for counter <= configuration.CommitCount {
		f := <- flows

		_, err = insert_query.Exec(f.TS, f.Proto, f.Src, f.Dst, f.Sport, f.Dport)

	        if err != nil {
    		    log.Print(err)
        	    break
	        }
		counter++
    	    }
	    insert_query.Close()

	    err = tx.Commit()
	    if err != nil {
		log.Print(err)
                time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
		break
	    }
	}
    }
}

func purgeDB(configuration Configuration) {
    connect_string := "host="+configuration.DBhost+" port="+configuration.DBport+" user="+configuration.DBuser+" password="+configuration.DBpassword+" dbname="+configuration.DBname+" sslmode=disable"
    
    ticker := time.NewTicker(time.Hour * PURGE_PERIOD)

    for t := range ticker.C {
	log.Print("Purging database at %v", t)
	db, err := sql.Open("postgres", connect_string)
	if err != nil {
    	    log.Print(err)
            break
	}

	db.Exec("delete from events where ts < (now() - interval '1 month')")
	println("purge")
	db.Close()
    }
}

