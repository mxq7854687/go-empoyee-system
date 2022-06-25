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
		log.Fatal("cannot load config: ", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("DB Connection [ Failed ]: ", err)
	}

	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	err = server.Start(config.ServerAdress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
