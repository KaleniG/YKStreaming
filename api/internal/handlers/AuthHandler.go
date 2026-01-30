package handlers

import (
	"log"
	"net/http"

	"ykstreaming_api/internal/db"
	"ykstreaming_api/internal/db/sqlc"
	"ykstreaming_api/internal/helpers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// INITIAL WEBSITE AUTHENTICATION CHECK
func Check(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			rememberToken, err := c.Cookie("remember_token")
			if err != nil || rememberToken == "" {
				if err == http.ErrNoCookie || rememberToken == "" {
					guestToken := session.Get("guest_token")
					if guestToken == nil {
						newGuestToken, err := helpers.GenerateRandomToken(32)
						if err != nil {
							log.Panic(err)
						}

						session.Set("guest_token", newGuestToken)
						if err := session.Save(); err != nil {
							log.Panic(err)
						}
					}

					c.JSON(http.StatusOK, gin.H{"user": false})
					return
				}
				log.Panic(err)
			}

			ctx := c.Request.Context()
			var rememberTokenText pgtype.Text
			if err := rememberTokenText.Scan(rememberToken); err != nil {
				log.Panic(err)
			}

			user, err := dbStore.Queries.GetUserDataByRememberToken(ctx, rememberTokenText)
			if err != nil {
				if err == pgx.ErrNoRows {
					c.SetCookie("remember_token", "", -1, "/", "localhost", false, true)
					c.JSON(http.StatusNotFound, gin.H{"error": "invalid remember me token, the token will be eliminated"})
					return
				}
				log.Panic(err)
			}

			session.Set("user_id", user.ID)
			if err := session.Save(); err != nil {
				log.Panic(err)
			}

			c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": user.Name, "email": user.Email}})
			return
		}

		ctx := c.Request.Context()
		user, err := dbStore.Queries.GetUserDataById(ctx, userID.(int32))
		if err != nil {
			if err == pgx.ErrNoRows {
				session.Delete("user_id")
				err = session.Save()
				if err != nil {
					log.Panic(err)
				}
				c.JSON(http.StatusNotFound, gin.H{"error": "invalid session user id, the user id will be eliminated"})
				return
			}
			log.Panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": user.Name, "email": user.Email}})
	}
}

// AUTHENTICATION MANAGEMENT
type loginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe *bool  `json:"remember_me" binding:"required"`
}

func Login(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login request"})
			return
		}

		ctx := c.Request.Context()
		user, err := dbStore.Queries.GetUserCredentialsByEmail(ctx, req.Email)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"param": "email", "error": "invalid email"})
				return
			}

			log.Panic(err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				c.JSON(http.StatusUnauthorized, gin.H{"param": "password", "error": "invalid password"})
				return
			}
			log.Panic(err)
		}

		if *req.RememberMe {
			rememberToken, err := helpers.GenerateRandomToken(32)
			if err != nil {
				log.Panic(err)
			}

			var rememberTokenText pgtype.Text
			if err := rememberTokenText.Scan(rememberToken); err != nil {
				log.Panic(err)
			}

			params := sqlc.UpdateUserRememberTokenParams{
				RememberToken: rememberTokenText,
				UserID:        user.ID,
			}

			err = dbStore.Queries.UpdateUserRememberToken(ctx, params)

			if err != nil {
				log.Panic(err)
			}

			c.SetCookie("remember_token", rememberToken, 2628000, "/", "localhost", false, true)
		} else {
			session := sessions.Default(c)
			session.Set("user_id", user.ID)
			err = session.Save()
			if err != nil {
				log.Panic(err)
			}
		}

		c.Status(http.StatusOK)
	}
}

type signupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Signup(dbStore *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signup request"})
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Panic(err)
		}

		ctx := c.Request.Context()
		userID, err := dbStore.Queries.AddUser(ctx, sqlc.AddUserParams{Name: req.Name, Email: req.Email, PasswordHash: string(passwordHash)})
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusConflict, gin.H{"param": "email", "error": "email conflict"})
				return
			}
			log.Panic(err)
		}

		session := sessions.Default(c)
		session.Set("user_id", userID)
		err = session.Save()
		if err != nil {
			log.Panic(err)
		}
		c.Status(http.StatusOK)
	}
}
