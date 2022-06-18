package db

import (
	"context"
	"database/sql"
	"example/employee/server/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load configuration with err: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("U cannot use the connection")
	}

	testQueries = New(testDB)

	code := m.Run()
	testQueries.DeleteAllRole(context.Background())
	os.Exit(code)
}
