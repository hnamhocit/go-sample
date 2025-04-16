package router

import (
	"sample/internal/config"
	"sample/internal/handlers"
	"sample/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Run() {
	// Load config
	dao, ctx := config.LoadConfig()

	// Initialize router
	router := gin.Default()

	// group
	api := router.Group("api")
	{
		auth := api.Group("auth")
		{
			authHandler := handlers.AuthHandler{Dao: dao, Ctx: ctx}
			auth.POST("login", authHandler.Login)
			auth.POST("register", authHandler.Register)
			auth.GET("refresh", middlewares.RefreshTokenMiddleware(dao, ctx), authHandler.Refresh)
			auth.GET("logout", middlewares.RefreshTokenMiddleware(dao, ctx), authHandler.Logout)
		}

		users := api.Group("users")
		{
			usersHandler := handlers.UserHandler{Dao: dao, Ctx: ctx}
			users.GET("", usersHandler.GetUsers)
			users.GET("me", middlewares.AccessTokenMiddleware(dao, ctx), usersHandler.GetUser)
			users.GET(":id", usersHandler.GetUser)
		}
	}

	router.Run()
}
