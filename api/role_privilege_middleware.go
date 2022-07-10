package api

import (
	db "example/employee/server/db/sqlc"
	"example/employee/server/service/role_service"
	"example/employee/server/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func privilegeMiddleware(roleService role_service.RoleService, privilege db.Privilege) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authPayLoadKey).(*token.Payload)

		authUser, err := roleService.Store.GetUser(ctx, authPayload.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		err = roleService.HasRolePriviledgeByRoleId(authUser.RoleID.Int64, privilege)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
		}

		ctx.Next()
	}
}
