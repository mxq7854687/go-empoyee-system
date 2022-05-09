package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateDepartmentRequest struct {
	departmentName string `json:"departmentName" binding:"required"`
}

func (server *Server) createDepartment(ctx *gin.Context) {
	var req CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	department, err := server.store.CreateDepartment(ctx, req.departmentName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, department)
}
