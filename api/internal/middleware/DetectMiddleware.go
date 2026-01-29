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
	if token != nil {
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
			// Try remember token
			rememberToken, err := c.Cookie("remember_token")
			if err != nil && err != http.ErrNoCookie {
				log.Print("Cookie read error:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if err == http.ErrNoCookie {
				// No cookie, fallback to guest
				if err := guestFallback(session); err != nil {
					log.Print(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			// Validate remember token
			var tokenText pgtype.Text
			if err := tokenText.Scan(rememberToken); err != nil {
				log.Print("Failed to scan remember token:", err)
				if err := guestFallback(session); err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			userID, err = dbStore.RQueries.GetUserIdByRememberToken(c.Request.Context(), tokenText)
			if err != nil {
				if err == pgx.ErrNoRows {
					c.SetCookie("remember_token", "", -1, "/", "localhost", false, true)
					log.Print("Invalid remember token cleared")
				} else {
					log.Print(err)
				}
				if err := guestFallback(session); err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			session.Set("user_id", userID)
			if err := session.Save(); err != nil {
				log.Print("Failed to save session:", err)
				if err := guestFallback(session); err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Next()
				return
			}

			c.Next()
			return
		}

		// Validate session user_id
		_, err := dbStore.RQueries.CheckUserById(c.Request.Context(), userID.(int32))
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Print("Invalid session user_id cleared")
				session.Delete("user_id")
				_ = session.Save()
			} else {
				log.Print(err)
			}
			if err := guestFallback(session); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		c.Next()
	}
}
