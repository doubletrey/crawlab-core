package middlewares

import (
	"github.com/doubletrey/crawlab-core/constants"
	"github.com/doubletrey/crawlab-core/controllers"
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/user"
	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	userSvc, _ := user.GetUserService()
	return func(c *gin.Context) {
		// token string
		tokenStr := c.GetHeader("Authorization")

		// validate token
		u, err := userSvc.CheckToken(tokenStr)
		if err != nil {
			// validation failed, return error response
			controllers.HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
			return
		}

		// set user in context
		c.Set(constants.UserContextKey, u)

		// validation success
		c.Next()
	}
}
