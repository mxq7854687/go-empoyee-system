package api

import (
	"database/sql"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type createEmployeeRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	JobID        int64  `json:"job_id" binding:"required"`
	Salary       int64  `json:"salary" binding:"required"`
	ManagerId    int64  `json:"manager_id" binding:"required"`
	DepartmentID int64  `json:"department_id" binding:"required"`
}

type employeeByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListEmployeesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listEmployees(ctx *gin.Context) {
	var req ListEmployeesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListEmployeesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	employees, err := server.store.ListEmployees(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, employees)
}

func (server *Server) getEmployeeById(ctx *gin.Context) {
	var req employeeByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	employee, err := server.store.GetEmployee(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

func (server *Server) createEmployee(ctx *gin.Context) {
	var req createEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEmployeeParams{
		FirstName:    util.GetSqlNullString(req.FirstName),
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  util.GetSqlNullString(req.PhoneNumber),
		HireDate:     time.Now(),
		JobID:        req.JobID,
		Salary:       req.Salary,
		ManagerID:    util.GetSqlNullInt64(req.ManagerId),
		DepartmentID: req.DepartmentID,
	}

	employee, err := server.store.CreateEmployee(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, employee)
}
