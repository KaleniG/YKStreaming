package handlers

import (
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	_ "golang.org/x/image/webp"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/db/sqlc"
	"ykstreaming_api/internal/helpers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// AUTHENTICATION MANAGEMENT
func Logout(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID != nil {
			session.Delete("user_id")
			err := session.Save()
			if err != nil {
				log.Panic(err)
			}
		}

		rememberToken, err := c.Cookie("remember_token")
		if err == nil {
			ctx := c.Request.Context()
			if userID != nil {
				err = dbStore.Queries.ResetUserRememberToken(ctx, userID.(int32))
			} else {
				var rememberTokenText pgtype.Text
				if err := rememberTokenText.Scan(rememberToken); err != nil {
					log.Panic(err)
				}
				err = dbStore.Queries.ResetRememberToken(ctx, rememberTokenText)
			}
			if err != nil {
				log.Panic(err)
			}

			c.SetCookie("remember_token", "", -1, "/http", "localhost", false, true)
		}

		c.Status(http.StatusOK)
	}
}

// STREAMS MANAGEMENT
func GetUserStreams(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			log.Panic(`session value "user_id" not set in auth protected route`)
		}

		ctx := c.Request.Context()
		streams, err := dbStore.Queries.GetUserStreams(ctx, userID.(int32))
		if err != nil {
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"streams": streams})
	}
}

type addStreamRequest struct {
	Name  string                `form:"name" binding:"required"`
	IsVOD *bool                 `form:"is_vod" binding:"required"`
	File  *multipart.FileHeader `form:"thumbnail"`
}

func AddStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addStreamRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream add request"})
			return
		}

		ctx := c.Request.Context()
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			log.Panic(`session value "user_id" not set in auth protected route`)
		}

		var isVodBool pgtype.Bool
		isVodBool.Scan(*req.IsVOD)

		var hasCustomThumbnailBool pgtype.Bool
		hasCustomThumbnailBool.Scan(req.File != nil)

		var newKey string
		var err error

		params := sqlc.AddStreamParams{
			UserID:             userID.(int32),
			Name:               req.Name,
			IsVod:              isVodBool,
			HasCustomThumbnail: hasCustomThumbnailBool,
		}

		for {
			newKey, err = helpers.GenerateRandomToken(64)
			if err != nil {
				log.Panic(err)
			}
			params.Key = newKey
			_, err = dbStore.Queries.AddStream(ctx, params)
			if err != nil {
				if err == pgx.ErrNoRows {
					continue
				}
				log.Panic(err)
			}
			break
		}

		if hasCustomThumbnailBool.Bool {
			customThumbnailsDir := helpers.GetEnvDir("CUSTOM_THUMBNAILS_DIR")

			f, err := req.File.Open()
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				log.Panic(err)
			}

			outputPath := customThumbnailsDir + "/" + newKey + ".jpg"
			out, err := os.Create(outputPath)
			if err != nil {
				log.Panic(err)
			}
			defer out.Close()

			err = jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
			if err != nil {
				log.Panic(err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"key": newKey})
	}
}

func RemoveStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamKey := c.Param("key")
		if strings.TrimSpace(streamKey) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream remove request"})
			return
		}

		ctx := c.Request.Context()
		data, err := dbStore.Queries.GetStreamRemovalDataByKey(ctx, streamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "streaming key not found"})
				return
			}
			log.Panic(err)
		}

		if data.IsActive.Bool {
			if data.IsVod.Bool {
				vodsDir := helpers.GetEnvDir("VODS_DIR")
				vodFLVFilepath := vodsDir + "/" + streamKey + ".flv"
				vodMP4Filepath := vodsDir + "/" + streamKey + ".mp4"

				err = requestStreamRecordingAction(streamKey, StopRecording)
				if err != nil {
					log.Panic(err)
				}

				if helpers.FileExists(vodFLVFilepath) {
					err = os.Remove(vodFLVFilepath)
					if err != nil {
						if !errors.Is(err, fs.ErrNotExist) {
							log.Panic(err)
						}
					}
				}

				if helpers.FileExists(vodMP4Filepath) {
					err = os.Remove(vodMP4Filepath)
					if err != nil {
						if !errors.Is(err, fs.ErrNotExist) {
							log.Panic(err)
						}
					}
				}
			}

			err = requestStreamStop(streamKey)
			if err != nil {
				log.Panic(err)
			}
		}

		err = dbStore.Tx(ctx, func(q *sqlc.Queries) error {
			err = q.RemoveStreamViewers(ctx, data.ID)
			if err != nil {
				return err
			}

			err = q.RemoveStream(ctx, data.ID)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		if data.HasCustomThumbnail.Bool {
			customThumbnailsDir := helpers.GetEnvDir("CUSTOM_THUMBNAILS_DIR")
			customThumbnailFilepath := customThumbnailsDir + "/" + streamKey + ".jpg"
			if helpers.FileExists(customThumbnailFilepath) {
				err = os.Remove(customThumbnailFilepath)
				if err != nil {
					if !errors.Is(err, fs.ErrNotExist) {
						log.Panic(err)
					}
				}
			} else {
				log.Print("the stream '" + streamKey + "' had a custom thumbnail but the file did not exists on deletion")
			}
		}

		liveThumbnailsDir := helpers.GetEnvDir("LIVE_THUMBNAILS_DIR")
		liveThumbnailFilepath := liveThumbnailsDir + "/" + streamKey + ".jpg"
		if helpers.FileExists(liveThumbnailFilepath) {
			err = os.Remove(liveThumbnailFilepath)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					log.Panic(err)
				}
			}
		} else {
			// NOT HANDLED SINCE LIVE THUMBNAILS ARE MANDATORY AND NOT OPTIONAL BUT TODO
		}

		liveThumbnailsLocksDir := helpers.GetEnvDir("LIVE_THUMBNAILS_LOCKS_DIR")
		liveThumbnailLockFilepath := liveThumbnailsLocksDir + "/" + streamKey + ".lock"
		if helpers.FileExists(liveThumbnailLockFilepath) {
			err = os.Remove(liveThumbnailLockFilepath)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					log.Panic(err)
				}
			}
		} else {
			// NOT HANDLED SINCE LIVE THUMBNAILS ARE MANDATORY AND NOT OPTIONAL BUT TODO
		}

		// JUST IN CASE THE STREAM FILES ARE NOT DELETED FOR SOME REASON
		// MAYBE NEED TO TRACK AND LOG THIS
		streamsDir := helpers.GetEnvDir("STREAMS_DIR")
		fileDirs, err := helpers.FindFilesContaining(streamsDir, streamKey)
		for _, fileDir := range fileDirs {
			err = os.Remove(fileDir)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					log.Panic(err)
				}
			}
		}

		c.Status(http.StatusOK)
	}
}

func StopStream(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamKey := c.Param("key")
		if strings.TrimSpace(streamKey) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream remove request"})
			return
		}

		ctx := c.Request.Context()
		data, err := dbStore.Queries.GetStreamStopDataByKey(ctx, streamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "streaming key not found"})
				return
			}
			log.Panic(err)
		}

		if data.IsActive.Bool {
			if data.IsVod.Bool {
				err = requestStreamRecordingAction(streamKey, StopRecording)
				if err != nil {
					log.Panic(err)
				}
			}

			err = requestStreamStop(streamKey)
			if err != nil {
				log.Panic(err)
			}

			err = dbStore.Queries.StopStream(ctx, data.ID)
			if err != nil {
				log.Panic(err)
			}
		}

		c.Status(http.StatusOK)
	}
}
