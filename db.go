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
    	    log.Fatal(err)
            time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
            continue
	}

	for {
	    counter := 0

	    tx, err := db.Begin()
	    if err != nil {
		log.Fatal(err)
                time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
		break
	    }

	    insert_query, err := tx.Prepare("insert into events (ts, proto, src, dst, sport, dport) values ($1, $2, $3, $4, $5, $6)")
	    defer insert_query.Close()

	    if err != nil {
		log.Fatal(err)
        	time.Sleep(10000 * time.Millisecond) //Wait 10 seconds before reconnect
        	break
    	    }

	    for counter <= configuration.CommitCount {
		f := <- flows

		_, err = insert_query.Exec(f.TS, f.Proto, f.Src, f.Dst, f.Sport, f.Dport)

	        if err != nil {
    		    log.Fatal(err)
        	    break
	        }
		counter++
    	    }
	    err = tx.Commit()
	    if err != nil {
		log.Fatal(err)
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
    	    log.Fatal(err)
            break
	}

	db.Exec("delete from events where ts < (now() - interval '1 month')")
	println("purge")
	db.Close()
    }
}

