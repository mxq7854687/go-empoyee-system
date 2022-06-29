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
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/departments", server.createDepartment)

	// router for job
	router.POST("/jobs", server.createJob)
	router.GET("/jobs/:id", server.getJob)
	router.GET("/jobs", server.listJobs)
	router.PUT("/jobs/:id", server.updateJob)
	router.DELETE("/jobs/:id", server.deleteJob)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
