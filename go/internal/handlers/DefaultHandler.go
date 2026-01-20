package handlers

import (
	"log"
	"net/http"
	"strings"
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/db/sqlc"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// STREAM INFORMATION
func GetStreams(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		streams, err := dbStore.Queries.GetPublicStreams(ctx)
		if err != nil {
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"streams": streams})
	}
}

func GetStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamKey := c.Param("key")
		if strings.TrimSpace(streamKey) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream status request"})
			return
		}

		ctx := c.Request.Context()
		stream, err := dbStore.Queries.GetStreamStatus(ctx, streamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "invalid stream key or inaccesible stream"})
				return
			}
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"stream": stream})
	}
}

// STREAM VIEWRSHIP MANAGEMENT
func ViewStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamKey := c.Param("key")
		if strings.TrimSpace(streamKey) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream view submit request"})
			return
		}

		ctx := c.Request.Context()
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID != nil {
			params := sqlc.ViewStreamAsUserParams{
				UserID:    userID.(int32),
				StreamKey: streamKey,
			}

			streamFound, err := dbStore.Queries.ViewStreamAsUser(ctx, params)
			if err != nil {
				log.Panic(err)
			}
			if !streamFound {
				log.Panic("stream key is invalid, user view not counted")
			}
		} else {
			guestToken := session.Get("guest_token")
			if guestToken != nil {
				params := sqlc.ViewStreamAsGuestParams{
					GuestToken: guestToken.(string),
					StreamKey:  streamKey,
				}

				streamFound, err := dbStore.Queries.ViewStreamAsGuest(ctx, params)
				if err != nil {
					log.Panic(err)
				}
				if !streamFound {
					log.Panic("stream key is invalid, user view not counted")
				}
			} else {
				log.Panic("user detection error")
			}
		}

		c.Status(http.StatusOK)
	}
}

func UnviewStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamKey := c.Param("key")
		if strings.TrimSpace(streamKey) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream unview submit request"})
			return
		}

		ctx := c.Request.Context()
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID != nil {
			params := sqlc.UnviewStreamAsUserParams{
				UserID:    userID.(int32),
				StreamKey: streamKey,
			}

			err := dbStore.Queries.UnviewStreamAsUser(ctx, params)
			if err != nil {
				log.Panic(err)
			}
		} else {
			guestToken := session.Get("guest_token")
			if guestToken != nil {
				params := sqlc.UnviewStreamAsGuestParams{
					GuestToken: guestToken.(string),
					StreamKey:  streamKey,
				}

				err := dbStore.Queries.UnviewStreamAsGuest(ctx, params)
				if err != nil {
					log.Panic(err)
				}
			} else {
				log.Panic("user detection error")
			}
		}

		c.Status(http.StatusOK)
	}
}
