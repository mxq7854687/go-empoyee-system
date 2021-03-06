package api

import (
	"database/sql"
	"errors"
	db "example/employee/server/db/sqlc"
	"example/employee/server/token"
	"example/employee/server/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
func (server *Server) login(ctx *gin.Context) {
	var req LoginRequest
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

	if user.Status != db.UserStatusActivated {
		err = fmt.Errorf("user not activated")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = util.ValidPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// TODO platform would offer different time of refresh token for user
	accessToken, err := server.tokenMaker.CreateToken(
		user.Email,
		token.Web,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LoginResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

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
