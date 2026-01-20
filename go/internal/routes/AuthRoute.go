package routes

import (
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup, dbStore *db.Store) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/check", handlers.Check(dbStore))
		authRouter.POST("/signup", handlers.Signup(dbStore))
		authRouter.POST("/login", handlers.Login(dbStore))
	}
}
