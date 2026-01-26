package testsetups

import (
	"log"
	"net/http"
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/helpers"
	"ykstreaming_api/internal/middleware"
	"ykstreaming_api/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRTMPRouter() (*gin.Engine, *db.Store) {
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

	authRouter := httpRoute.Group("/auth")
	{
		authRouter.POST("/check", handlers.Check(dbStore))
		authRouter.POST("/signup", handlers.Signup(dbStore))
	}

	userRouter := httpRoute.Group("/user")
	userRouter.Use(middleware.Auth(dbStore))
	userStreamsRouter := userRouter.Group("/streams")
	{
		userStreamsRouter.POST("/add", handlers.AddStream(dbStore))
		userStreamsRouter.POST("/stop/:key", handlers.StopStream(dbStore))
	}
	routes.RTMPRoute(router, dbStore)

	return router, dbStore
}
