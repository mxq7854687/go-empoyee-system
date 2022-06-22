package api

import (
	db "example/employee/server/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/departments", server.createDepartment)

	router.POST("/auth/activate", server.activateUser)

	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
