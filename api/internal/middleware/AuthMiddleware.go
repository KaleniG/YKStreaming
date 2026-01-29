package middleware

import (
	"log"
	"net/http"

	"ykstreaming_api/internal/db"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func Auth(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			rememberToken, err := c.Cookie("remember_token")
			if err != nil {
				if err == http.ErrNoCookie {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				log.Panic(err)
			}

			var rememberTokenText pgtype.Text
			if err := rememberTokenText.Scan(rememberToken); err != nil {
				log.Panic(err)
			}

			ctx := c.Request.Context()
			userID, err = dbStore.RQueries.GetUserIdByRememberToken(ctx, rememberTokenText)
			if err != nil {
				if err == pgx.ErrNoRows {
					c.SetCookie("remember_token", "", -1, "/", "localhost", false, true)
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "invalid remember me token, the token will be eliminated"})
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}

			session.Set("user_id", userID)
			if err := session.Save(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}

			c.Next()
			return
		}

		ctx := c.Request.Context()
		_, err := dbStore.RQueries.CheckUserById(ctx, userID.(int32))
		if err != nil {
			if err == pgx.ErrNoRows {
				session.Delete("user_id")
				session.Save()
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "invalid session user id, the user id will be eliminated"})
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Next()
	}
}
