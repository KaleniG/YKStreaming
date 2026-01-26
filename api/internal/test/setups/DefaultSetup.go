package testsetups

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/middleware"
	"ykstreaming_api/internal/routes"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRouter() (*gin.Engine, *db.Store) {
	engine, router, dbStore := setupHTTPRouter()

	routes.AuthRoute(router, dbStore)
	routes.DefaultRoute(router, dbStore)

	// Partial UserRoute
	userRouter := router.Group("/user")
	userRouter.Use(middleware.Auth(dbStore))
	{
		userRouter.POST("/logout", handlers.Logout(dbStore))
	}

	userStreamsRouter := userRouter.Group("/streams")
	{
		userStreamsRouter.POST("/add", handlers.AddStream(dbStore))
		userStreamsRouter.POST("/remove/:key", handlers.RemoveStream(dbStore))
	}

	routes.RTMPRoute(engine, dbStore)

	return engine, dbStore
}
