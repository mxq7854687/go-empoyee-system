package api

import (
	"database/sql"
	db "example/employee/server/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateJobRequest struct {
	JobTitle  string        `form:"job_title" json:"job_title" binding:"required"`
	MinSalary sql.NullInt64 `form:"min_salary" json:"min_salary" binding:"required"`
	MaxSalary sql.NullInt64 `from:"max_salary" json:"max_salary" binding:"required"`
}

func (server *Server) createJob(ctx *gin.Context) {
	var req CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateJobParams{
		JobTitle:  req.JobTitle,
		MinSalary: req.MinSalary,
		MaxSalary: req.MaxSalary,
	}

	job, err := server.store.CreateJob(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}

type GetJobRequest struct {
	JobID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getJob(ctx *gin.Context) {
	var req GetJobRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	job, err := server.store.GetJob(ctx, req.JobID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}

type ListJobsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listJobs(ctx *gin.Context) {
	var req ListJobsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListJobsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	jobs, err := server.store.ListJobs(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

type UpdateJobRequestUri struct {
	JobID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateJobRequestJson struct {
	JobTitle  string        `form:"job_title" json:"job_title" binding:"required"`
	MinSalary sql.NullInt64 `form:"min_salary" json:"min_salary" binding:"required"`
	MaxSalary sql.NullInt64 `from:"max_salary" json:"max_salary" binding:"required"`
}

func (server *Server) updateJob(ctx *gin.Context) {
	var uri UpdateJobRequestUri
	var req UpdateJobRequestJson
	if errInUri := ctx.ShouldBindUri(&uri); errInUri != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errInUri))
		return
	}

	if errInJson := ctx.ShouldBindJSON(&req); errInJson != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errInJson))
		return
	}

	arg := db.UpdateJobParams{
		JobID:     uri.JobID,
		JobTitle:  req.JobTitle,
		MinSalary: req.MinSalary,
		MaxSalary: req.MaxSalary,
	}

	job, err := server.store.UpdateJob(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}

type DeleteJobRequest struct {
	JobID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteJob(ctx *gin.Context) {
	var req DeleteJobRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err1 := server.store.GetJob(ctx, req.JobID)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err1))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err1))
		return
	}

	err2 := server.store.DeleteJob(ctx, req.JobID)

	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err2))
	}

	ctx.JSON(http.StatusOK, "success")
}
