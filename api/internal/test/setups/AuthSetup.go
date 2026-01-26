package testsetups

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/middleware"
	"ykstreaming_api/internal/routes"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter() (*gin.Engine, *db.Store) {
	engine, router, dbStore := setupHTTPRouter()

	routes.AuthRoute(router, dbStore)

	// Partial UserRoute
	userRouter := router.Group("/user")
	userRouter.Use(middleware.Auth(dbStore))
	{
		userRouter.POST("/logout", handlers.Logout(dbStore))
	}

	return engine, dbStore
}
