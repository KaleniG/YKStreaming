package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbStore := db.OpenDefault()
	defer dbStore.Close()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is missing")
	}

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	routes.HTTPRoute(router, dbStore)
	routes.RTMPRoute(router, dbStore)

	err = router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal(err)
	}
}
