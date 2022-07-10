package api

import (
	db "example/employee/server/db/sqlc"
	"example/employee/server/service/role_service"
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

func NewServer(config util.Config, store db.Store, roleService role_service.RoleService) (*Server, error) {

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

	router.POST("/auth/activate", server.activateUser)
	router.POST("/auth/login", server.login)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/departments",
		privilegeMiddleware(roleService, db.PrivilegeCreateAndUpdateJobs),
		server.createDepartment)
	// router for job
	authRoutes.POST("/jobs",
		privilegeMiddleware(roleService, db.PrivilegeCreateAndUpdateJobs),
		server.createJob)
	authRoutes.GET("/jobs/:id", server.getJob)
	authRoutes.GET("/jobs", server.listJobs)
	authRoutes.PUT("/jobs/:id",
		privilegeMiddleware(roleService, db.PrivilegeCreateAndUpdateJobs),
		server.updateJob)
	authRoutes.DELETE("/jobs/:id",
		privilegeMiddleware(roleService, db.PrivilegeDeleteJobs),
		server.deleteJob)

	server.router = router
	return server, nil
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
