package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {

	var err error
    connStr := "user=yosa dbname=jewelerysitedb sslmode=disable password=yosa"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    err = db.Ping()
    if err != nil {
        log.Fatal("Could not ping the database:", err)
    }
    fmt.Println("Connected to the database")
}
func GetDB() *sql.DB {
	return db
}

