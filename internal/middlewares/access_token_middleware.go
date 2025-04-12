package middlewares

import (
	"net/http"
	"sample/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation of access token middleware\
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"code": 0,
					"msg":  "Unauthorized",
				})
			return
		}

		tokenString := strings.Split(authorization, " ")[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"code": 0,
					"msg":  "Invalid token",
				})
			return
		}

		claims, err := utils.VerifyToken(tokenString, "JWT_ACCESS_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"code": 0,
					"msg":  err.Error(),
				})
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"code": 0,
					"msg":  err.Error(),
				})
			return
		}

		// save user id in context
		c.Set("user_id", sub)

		c.Next()
	}
}
