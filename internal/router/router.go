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
			userHandler := handlers.UserHandler{Dao: dao, Ctx: ctx}
			users.GET("", userHandler.GetUsers)
			users.GET("me", middlewares.AccessTokenMiddleware(dao, ctx), userHandler.GetUser)
			users.GET(":id", userHandler.GetUser)
		}

		media := api.Group("media")
		{
			mediaHandler := handlers.MediaHandler{Dao: dao, Ctx: ctx}
			media.POST("upload", middlewares.AccessTokenMiddleware(dao, ctx), mediaHandler.Upload)
			media.POST("uploads", middlewares.AccessTokenMiddleware(dao, ctx), mediaHandler.Uploads)
		}
	}

	router.Static("/assets", "./assets")
	router.Run()
}
