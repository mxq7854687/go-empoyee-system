package api

import (
	"database/sql"
	"errors"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ActivateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ActivateUserResponse struct {
	Email     string        `json:"email"`
	Status    db.UserStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (server *Server) activateUser(ctx *gin.Context) {
	var req ActivateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Status != db.UserStatusPending {
		err := errors.New("Only pending user can be activated")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ActivateUserParams{
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}
	err = server.store.ActivateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	user, err = server.store.GetUser(ctx, req.Email)
	resp := ActivateUserResponse{
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
