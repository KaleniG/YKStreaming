package routes

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RTMPRoute(router *gin.Engine, dbStore *db.Store) {
	rtmpRoute := router.Group("/rtmp")
	{
		rtmpRoute.POST("/on-publish", handlers.OnStreamPublish(dbStore))
		rtmpRoute.POST("/on-publish-done", handlers.OnStreamPublishDone(dbStore))
		rtmpRoute.POST("/on-update", handlers.OnStreamUpdate(dbStore))
		rtmpRoute.POST("/on-record-done", handlers.OnStreamRecordDone(dbStore))
	}
}
