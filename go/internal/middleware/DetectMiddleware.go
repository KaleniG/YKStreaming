package middleware

import (
	"log"
	"net/http"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/helpers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func guestFallback(session sessions.Session) error {
	token := session.Get("guest_token")
	if token == nil {
		return nil
	}

	newGuestToken, err := helpers.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	session.Set("guest_token", newGuestToken)
	return session.Save()
}

func Detect(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			rememberToken, err := c.Cookie("remember_token")
			if err != nil {
				log.Print(err)
				if err := guestFallback(session); err != nil {
					log.Print(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			ctx := c.Request.Context()
			var rememberTokenText pgtype.Text
			if err := rememberTokenText.Scan(rememberToken); err != nil {
				log.Panic(err)
			}
			userID, err = dbStore.Queries.GetUserIdByRememberToken(ctx, rememberTokenText)
			if err != nil {
				if err == pgx.ErrNoRows {
					c.SetCookie("remember_token", "", -1, "/", "localhost", false, true)
					log.Print("invalid remember me token, the token will be eliminated")
				} else {
					log.Print(err)
				}
				if err := guestFallback(session); err != nil {
					log.Print(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			session.Set("user_id", userID)
			if err := session.Save(); err != nil {
				log.Print(err)
				if err := guestFallback(session); err != nil {
					log.Print(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			c.Next()
			return
		}

		ctx := c.Request.Context()
		_, err := dbStore.Queries.CheckUserById(ctx, userID.(int32))
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Print("invalid session user_id, the user_id will be eliminated")
				session.Delete("user_id")
				if err := session.Save(); err != nil {
					log.Print(err)
				}
			} else {
				log.Print(err)
			}
			if err := guestFallback(session); err != nil {
				log.Print(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Next()
			return
		}
		c.Next()
	}
}
