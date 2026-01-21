package routes

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DefaultRoute(router *gin.RouterGroup, dbStore *db.Store) {
	router.POST("/get-streams", handlers.GetStreams(dbStore))

	stream := router.Group("/stream")
	{
		stream.POST("/:key", handlers.GetStream(dbStore))

		streamWithUser := stream.Group("/", middleware.Detect(dbStore))
		{
			streamWithUser.POST("/view/:key", handlers.ViewStream(dbStore))
			streamWithUser.POST("/unview/:key", handlers.UnviewStream(dbStore))
		}
	}
}
