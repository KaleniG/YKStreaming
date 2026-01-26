package routes

import (
	"log"
	"net/http"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/helpers"
	"ykstreaming_api/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func HTTPRoute(router *gin.Engine, dbStore *db.Store) {
	sessionAuthKey, err := helpers.GetEnvDir("SESSION_AUTH_KEY")
	if err != nil {
		log.Fatal(err)
	}

	/*
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}))
	*/

	sessionStore := cookie.NewStore([]byte(sessionAuthKey))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false, // true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	httpRoute := router.Group("/http")

	httpRoute.Use(sessions.Sessions("yksession", sessionStore))
	httpRoute.Use(middleware.CORS())

	{
		httpRoute.OPTIONS("/*path", handlers.OptionsCORSHandler())
		DefaultRoute(httpRoute, dbStore)
		AuthRoute(httpRoute, dbStore)
		UserRoute(httpRoute, dbStore)
	}
}
