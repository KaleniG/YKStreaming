package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Router  *gin.Engine
	Cookies []*http.Cookie
}

func NewClient(router *gin.Engine) *Client {
	return &Client{Router: router}
}

func (c *Client) Do(req *http.Request) *httptest.ResponseRecorder {
	// Add existing cookies to request
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)

	// Merge cookies instead of replacing
	newCookies := w.Result().Cookies()
	for _, newCookie := range newCookies {
		found := false
		for i, oldCookie := range c.Cookies {
			if oldCookie.Name == newCookie.Name {
				// replace old cookie with new one
				c.Cookies[i] = newCookie
				found = true
				break
			}
		}
		if !found {
			c.Cookies = append(c.Cookies, newCookie)
		}
	}

	return w
}
