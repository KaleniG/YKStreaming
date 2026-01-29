package main

import (
	"context"
	"log"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"ykstreaming_api/internal/db"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbStore := db.OpenDefault()
	defer dbStore.Close()

	ctx := context.Background()
	var idle bool = false

	log.Print("[YK] WatchDB job started")
	for {
		if idle {
			log.Print("[YK] Idle")
			time.Sleep(10 * time.Second)
		} else {
			log.Print("[YK] Not Idle")
			time.Sleep(5 * time.Second)
		}

		supervisedStreams, err := dbStore.RQueries.GetViewerActiveStreams(ctx)
		if err != nil {
			log.Print(err)
			idle = true
			continue
		}
		log.Print("[YK] Apprehended " + strconv.Itoa(len(supervisedStreams)) + " viewer active streams")
		if len(supervisedStreams) == 0 {
			idle = true
			continue
		}
		idle = false

		for _, streamID := range supervisedStreams {
			err := dbStore.WQueries.UpdateStreamLiveViewersCountByID(ctx, streamID)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}
