package handlers

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type rtmpPublishRequest struct {
	StreamKey string `form:"name" binding:"required"`
	Addr      string `form:"addr"`
	App       string `form:"app"`
	Tcurl     string `form:"tcurl"`
	PageUrl   string `form:"pageUrl"`
}

func OnStreamPublish(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rtmpPublishRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid on stream publish nginx rtmp module request"})
			return
		}

		ctx := c.Request.Context()

		_, err := dbStore.Queries.CheckStreamExistsByKey(ctx, req.StreamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "the stream does not exist"})
				return
			}
			log.Panic(err)
		}

		streamStarted, err := dbStore.Queries.StartStream(ctx, req.StreamKey)
		if err != nil {
			log.Panic(err)
		}
		if !streamStarted {
			c.JSON(http.StatusForbidden, gin.H{"error": "the stream has already started or ended"})
		}

		c.Status(http.StatusOK)
	}
}

func OnStreamPublishDone(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rtmpPublishRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid on stream publish done nginx rtmp module request"})
			return
		}

		ctx := c.Request.Context()
		_, err := dbStore.Queries.CheckStreamExistsByKey(ctx, req.StreamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "the stream does not exist"})
				return
			}
			log.Panic(err)
		}

		isVOD, err := dbStore.Queries.EndStream(ctx, req.StreamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusForbidden, gin.H{"error": "the stream has already ended or did not start yet"})
				return
			}
			log.Panic(err)
		}
		if isVOD.Bool {
			err = requestStreamRecordingAction(req.StreamKey, StopRecording)
			if err != nil {
				log.Panic(err)
			}
		}

		c.Status(http.StatusOK)
	}
}

type rtmpUpdateRequest struct {
	StreamKey string `form:"name" binding:"required"`
	Addr      string `form:"addr"`
	App       string `form:"app"`
	Tcurl     string `form:"tcurl"`
	PageUrl   string `form:"pageUrl"`
	Time      *int64 `form:"time"`
	Timestamp string `form:"timestamp"`
}

func OnStreamUpdate(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rtmpUpdateRequest
		if err := c.ShouldBind(&req); err != nil {
			log.Panic(req)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid on stream update nginx rtmp module request"})
			return
		}

		if req.Time != nil && *req.Time == 0 {
			c.Status(http.StatusOK)
			return
		}

		ctx := c.Request.Context()
		streamStatus, err := dbStore.Queries.GetStreamStatus(ctx, req.StreamKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "the stream does not exist"})
				return
			}
			log.Panic(err)
		}

		if !streamStatus.IsLive.Bool {
			c.Status(http.StatusForbidden)
			return
		}

		vodsDir, err := helpers.GetEnvDir("VODS_DIR")
		if err != nil {
			log.Panic(err)
		}
		vodFLVFilepath := vodsDir + "/" + req.StreamKey + ".flv"
		if streamStatus.IsVod.Bool && !helpers.FileExists(vodFLVFilepath) {
			err = requestStreamRecordingAction(req.StreamKey, StartRecording)
			if err != nil {
				log.Panic(err)
			}
		}

		liveThumbnailsLocksDir, err := helpers.GetEnvDir("LIVE_THUMBNAILS_LOCKS_DIR")
		if err != nil {
			log.Panic(err)
		}
		lockFilepath := liveThumbnailsLocksDir + "/" + req.StreamKey + ".lock"

		info, err := os.Stat(lockFilepath)
		if err == nil {
			if time.Since(info.ModTime()) < 60*time.Second {
				c.Status(http.StatusOK)
				return
			}
		} else if !os.IsNotExist(err) {
			log.Panic(err)
		}

		f, err := os.OpenFile(lockFilepath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		now := time.Now()
		if err := os.Chtimes(lockFilepath, now, now); err != nil {
			log.Panic(err)
		}

		liveThumbnailsDir, err := helpers.GetEnvDir("LIVE_THUMBNAILS_DIR")
		if err != nil {
			log.Panic(err)
		}
		thumbnailFilepath := liveThumbnailsDir + "/" + req.StreamKey + ".jpg"
		rtmpStreamURL := "rtmp://localhost/live/" + req.StreamKey
		cmd := exec.Command(
			"sudo",
			"ffmpeg", "-y",
			"-i", rtmpStreamURL,
			"-frames:v", "1",
			"-q:v", "3",
			thumbnailFilepath,
		)
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Panic(err)
		}

		c.Status(http.StatusOK)
	}
}

type rtmpRecordEndRequest struct {
	Recorder string `form:"recorder"`
	Path     string `form:"path" binding:"required"`
}

func OnStreamRecordDone(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rtmpRecordEndRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid on stream record done nginx rtmp module request"})
			return
		}

		ext := filepath.Ext(req.Path)
		if ext != "" {
			req.Path = req.Path[:len(req.Path)-len(ext)]
		}

		flvFilepath := req.Path + ".flv"
		mp4Filepath := req.Path + ".mp4"

		if helpers.FileExists(flvFilepath) {
			cmd := exec.Command(
				"ffmpeg", "-y",
				"-i", flvFilepath,
				"-c", "copy",
				mp4Filepath,
			)
			//cmd.Stderr = os.Stderr
			//cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				cmd := exec.Command(
					"ffmpeg", "-y",
					"-i", flvFilepath,
					"-c:v", "libx264",
					"-c:a", "aac",
					mp4Filepath,
				)
				//cmd.Stderr = os.Stderr
				//cmd.Stdout = os.Stdout

				err := cmd.Run()
				if err != nil {
					log.Panic(err)
				}
			}
		}

		c.Status(http.StatusOK)
	}
}
