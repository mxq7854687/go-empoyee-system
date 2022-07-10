package main

import (
	"context"
	"database/sql"
	"example/employee/server/api"
	db "example/employee/server/db/sqlc"
	"example/employee/server/service/role_service"
	"example/employee/server/util"
	"log"

	_ "github.com/lib/pq"
)

var roleService role_service.RoleService

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
	roleService := role_service.NewRoleService(store, context.Background())
	roleService.InitRole()

	server, err := api.NewServer(config, store, *roleService)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	err = server.Start(config.ServerAdress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
