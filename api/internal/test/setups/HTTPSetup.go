package testsetups

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

func setupHTTPRouter() (*gin.Engine, *gin.RouterGroup, *db.Store) {
	router, dbStore := setupCoreRouter()

	sessionAuthKey, err := helpers.GetEnvDir("SESSION_AUTH_KEY")
	if err != nil {
		log.Fatal(err)
	}

	sessionStore := cookie.NewStore([]byte(sessionAuthKey))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	httpRoute := router.Group("/http")
	httpRoute.Use(sessions.Sessions("yksession", sessionStore))
	httpRoute.Use(middleware.CORS())

	{
		httpRoute.OPTIONS("/*path", handlers.OptionsCORSHandler())
		httpRoute.GET("/session/:name", sessionToolHandler())
	}

	return router, httpRoute, dbStore
}
