package role_service

import (
	"context"
	"database/sql"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDB *sql.DB
var roleService *RoleService

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load configuration with err: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("U cannot use the connection")
	}
	store := db.NewStore(testDB)
	roleService = NewRoleService(store, context.Background())

	code := m.Run()
	os.Exit(code)
}
