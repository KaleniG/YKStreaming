package test

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/handlers"
	"ykstreaming_api/internal/helpers"
	"ykstreaming_api/internal/middleware"
	"ykstreaming_api/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func sessionToolHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("name")

		session := sessions.Default(c)
		userID := session.Get(val)
		if userID == nil {
			c.String(200, "")
			return
		}
		switch v := userID.(type) {
		case string:
			c.String(200, v)
		case int:
			c.String(200, strconv.Itoa(v))
		case int32:
			c.String(200, strconv.FormatInt(int64(v), 10))
		case int64:
			c.String(200, strconv.FormatInt(v, 10))
		default:
			c.String(500, "unexpected user_id type")
		}
	}
}

func SetupAuthRouter() (*gin.Engine, *db.Store) {
	err := godotenv.Load("/var/www/html/api/.env")
	if err != nil {
		log.Fatal(err)
	}

	dbStore := db.Open()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is missing")
	}

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	sessionAuthKey := helpers.GetEnvDir("SESSION_AUTH_KEY")

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
		routes.AuthRoute(httpRoute, dbStore)
	}

	userRouter := httpRoute.Group("/user")
	userRouter.Use(middleware.Auth(dbStore))
	{
		userRouter.POST("/logout", handlers.Logout(dbStore))
	}

	return router, dbStore
}
