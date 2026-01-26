package testsetups

import (
	"context"
	"log"
	"os"
	"strconv"
	"ykstreaming_api/internal/db"

	"github.com/gin-contrib/sessions"
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

func clearDatabase(dbStore *db.Store) {
	ctx := context.Background()
	dbStore.Queries.RemoveAllViews(ctx)
	dbStore.Queries.RemoveAllStreams(ctx)
	dbStore.Queries.RemoveAllUsers(ctx)
}

func setupCoreRouter() (*gin.Engine, *db.Store) {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}

	dbStore := db.OpenTest()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is missing")
	}

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	clearDatabase(dbStore)

	return router, dbStore
}
