package api

import (
	db "example/employee/server/db/sqlc"
	"example/employee/server/token"
	"example/employee/server/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Not able to create token: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()

	router.POST("/auth/login", server.login)
	router.POST("/auth/activate", server.activateUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/departments", server.createDepartment)
	// router for job
	authRoutes.POST("/jobs", server.createJob)
	authRoutes.GET("/jobs/:id", server.getJob)
	authRoutes.GET("/jobs", server.listJobs)
	authRoutes.PUT("/jobs/:id", server.updateJob)
	authRoutes.DELETE("/jobs/:id", server.deleteJob)

	server.router = router
	return server, nil
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
