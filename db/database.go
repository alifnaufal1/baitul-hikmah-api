package db

import (
	"blog-api/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
	dbConfig, err := config.LoadConfig()
  if err != nil {
    log.Fatalln("error while loading .env file", err.Error())   
  }

  dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.DBName)

  DB, err = sql.Open("pgx", dsn)
  if err != nil {
    log.Fatalln("error while opening database", err.Error())
  }
  
  err = DB.Ping()
  if err != nil {
    log.Fatalln("error while pinging database", err.Error())
  }

  fmt.Println("success connect to db")
}

// migrate:
// migrate create -ext sql -dir db/migrations {nama-migrate}
// migrate -database "postgres://postgres:{pw}@localhost:5432/free-market-test?sslmode=disable" -path db/migrations up
