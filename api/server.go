package api

import (
	"context"
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

	roleService := role_service.NewRoleService(store, context.Background())
	roleService.InitRole()

	router := gin.Default()

	router.POST("/auth/activate", server.activateUser)
	router.POST("/auth/login", server.login)

	router.POST("/departments", server.createDepartment,
		privilegeMiddleware(context.Background(), *roleService, db.PrivilegeCreateAndUpdateDepartments))
	// router for job
	router.POST("/jobs", server.createJob).Use(
		privilegeMiddleware(context.Background(), *roleService, db.PrivilegeCreateAndUpdateJobs))
	router.GET("/jobs/:id", server.getJob)
	router.GET("/jobs", server.listJobs)
	router.PUT("/jobs/:id", server.updateJob,
		privilegeMiddleware(context.Background(), *roleService, db.PrivilegeCreateAndUpdateJobs))
	router.DELETE("/jobs/:id", server.deleteJob,
		privilegeMiddleware(context.Background(), *roleService, db.PrivilegeDeleteJobs))

	server.router = router
	return server, nil
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
