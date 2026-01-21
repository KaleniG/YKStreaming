package routes

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, dbStore *db.Store) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.Auth(dbStore))
	{
		userRouter.POST("/logout", handlers.Logout(dbStore))
	}

	userStreamsRouter := userRouter.Group("/streams")
	{
		userStreamsRouter.POST("/", handlers.GetUserStreams(dbStore))
		userStreamsRouter.POST("/add", handlers.AddStream(dbStore))
		userStreamsRouter.POST("/stop/:key", handlers.StopStream(dbStore))
		userStreamsRouter.POST("/remove/:key", handlers.RemoveStream(dbStore))
	}
}
