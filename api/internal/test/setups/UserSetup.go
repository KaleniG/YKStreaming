package testsetups

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/routes"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter() (*gin.Engine, *db.Store) {
	engine, router, dbStore := setupHTTPRouter()

	routes.AuthRoute(router, dbStore)
	routes.UserRoute(router, dbStore)
	routes.RTMPRoute(engine, dbStore)

	return engine, dbStore
}
