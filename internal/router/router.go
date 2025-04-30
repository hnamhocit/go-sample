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
	baseHandler := handlers.BaseHandler{Dao: dao, Ctx: ctx}
	api := router.Group("api")
	{
		auth := api.Group("auth")
		{
			authHandler := handlers.AuthHandler{BaseHandler: baseHandler}
			auth.POST("login", authHandler.Login)
			auth.POST("register", authHandler.Register)
			auth.GET("refresh", middlewares.RefreshTokenMiddleware(dao, ctx), authHandler.Refresh)
			auth.GET("logout", middlewares.RefreshTokenMiddleware(dao, ctx), authHandler.Logout)
		}

		users := api.Group("users")
		{
			userHandler := handlers.UserHandler{BaseHandler: baseHandler}
			users.GET("", userHandler.GetUsers)
			users.GET("me", middlewares.AccessTokenMiddleware(dao, ctx), userHandler.GetMe)
			users.GET(":id", userHandler.GetUser)
		}

		media := api.Group("media")
		{
			mediaHandler := handlers.MediaHandler{BaseHandler: baseHandler}
			media.POST("upload", middlewares.AccessTokenMiddleware(dao, ctx), mediaHandler.Upload)
			media.POST("uploads", middlewares.AccessTokenMiddleware(dao, ctx), mediaHandler.Uploads)
		}

		posts := api.Group("posts")
		{
			postHandler := handlers.PostHandler{BaseHandler: baseHandler}
			posts.GET("", postHandler.GetPosts)
			posts.POST("", middlewares.AccessTokenMiddleware(dao, ctx), postHandler.CreatePost)
			posts.PUT(":id", middlewares.AccessTokenMiddleware(dao, ctx), postHandler.UpdatePost)
			posts.DELETE(":id", middlewares.AccessTokenMiddleware(dao, ctx), postHandler.DeletePost)
			posts.GET(":id", postHandler.GetPost)
		}
	}

	router.Static("/assets", "./assets")
	router.Run()
}
