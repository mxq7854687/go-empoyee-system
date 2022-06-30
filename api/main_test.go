package api

import (
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
