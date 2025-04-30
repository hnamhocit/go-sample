package middlewares

import (
	"context"
	"net/http"
	"sample/internal/database"
	"sample/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshTokenMiddleware(dao *database.Queries, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusOK,
				gin.H{
					"code": 0,
					"msg":  "Unauthorized",
				})
			return
		}

		tokenString := strings.Split(authorization, " ")[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusOK,
				gin.H{
					"code": 0,
					"msg":  "Invalid token",
				})
			return
		}

		claims, err := utils.VerifyToken(tokenString, "JWT_REFRESH_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK,
				gin.H{
					"code": 0,
					"msg":  "Verify token error: " + err.Error(),
				})
			return
		}

		tokenVersion, err := dao.GetUserTokenVersion(ctx, claims.Sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK,
				gin.H{
					"code": 0,
					"msg":  err.Error(),
				})
			return
		}

		if tokenVersion != claims.TokenVersion {
			c.AbortWithStatusJSON(http.StatusOK,
				gin.H{
					"code": 0,
					"msg":  "Token version mismatch",
				})
			return
		}

		c.Set("user_id", claims.Sub)
		c.Set("refresh_token", tokenString)

		c.Next()
	}
}
