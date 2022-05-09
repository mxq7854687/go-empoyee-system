package main

import (
	"database/sql"
	"example/employee/server/api"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("Cannot load configuration with err: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAdress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
