package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/sebastian-nunez/golang-store-api/cmd/api"
	"github.com/sebastian-nunez/golang-store-api/config"
	"github.com/sebastian-nunez/golang-store-api/db"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal("DB: unable to connect to the database. ", err)
	}

	initStorage(db)

	server := api.NewServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal("Server: unable to run. ", err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("DB: unable to ping! ", err)
	}

	log.Println("DB: successfully connected!")
}
