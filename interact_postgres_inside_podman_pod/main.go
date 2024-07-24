package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	//make a deadline before a Postgres container started
	time.Sleep(3 * time.Second)
	connStr := "user=postgres password=secret port=5432 host=localhost dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("postgres error: ", err)
	}
	var version string
	err = db.QueryRow("select version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)

	defer db.Close()
}
