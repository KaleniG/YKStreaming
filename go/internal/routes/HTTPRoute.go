package routes

import (
	"net/http"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/helpers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func HTTPRoute(router *gin.Engine, dbStore *db.Store) {
	sessionAuthKey := helpers.GetEnvDir("SESSION_AUTH_KEY")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	cookieStore := cookie.NewStore([]byte(sessionAuthKey))
	cookieStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false, // true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	httpRoute := router.Group("/http")

	httpRoute.Use(sessions.Sessions("yksession", cookieStore))
	//httpRoute.Use(middleware.CORS())

	{
		/*
			httpRoute.OPTIONS("/*path", func(c *gin.Context) {
				c.Status(204)
			})
		*/
		DefaultRoute(httpRoute, dbStore)
		AuthRoute(httpRoute, dbStore)
		UserRoute(httpRoute, dbStore)
	}
}
