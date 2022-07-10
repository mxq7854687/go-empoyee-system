package api

import (
	"context"
	db "example/employee/server/db/sqlc"
	"example/employee/server/service/role_service"
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

	roleService := role_service.NewRoleService(store, context.Background())
	server, err := NewServer(config, store, *roleService)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
