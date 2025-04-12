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
			authRepo := handlers.AuthRepo{Dao: dao, Ctx: ctx}
			auth.POST("login", authRepo.Login)
			auth.POST("register", authRepo.Register)
			auth.GET("refresh", middlewares.RefreshTokenMiddleware(), authRepo.Refresh)
			auth.GET("logout", middlewares.AccessTokenMiddleware(), authRepo.Logout)
		}
	}

	router.Run()
}
