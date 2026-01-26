package testhelpers

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Client struct {
	Router  *gin.Engine
	Cookies []*http.Cookie
}

func NewClient(router *gin.Engine) *Client {
	return &Client{Router: router}
}

func (client *Client) Do(req *http.Request) *httptest.ResponseRecorder {
	// Add existing cookies to request
	for _, cookie := range client.Cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	client.Router.ServeHTTP(w, req)

	// Merge cookies instead of replacing
	newCookies := w.Result().Cookies()
	for _, newCookie := range newCookies {
		found := false
		for i, oldCookie := range client.Cookies {
			if oldCookie.Name == newCookie.Name {
				// replace old cookie with new one
				client.Cookies[i] = newCookie
				found = true
				break
			}
		}
		if !found {
			client.Cookies = append(client.Cookies, newCookie)
		}
	}

	return w
}

func (client *Client) GetCookieKey(name string) string {
	resCookies := client.Cookies
	for _, c := range resCookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func (client *Client) GetSessionKey(name string) string {
	req, _ := http.NewRequest("GET", "/http/session/"+name, nil)
	res := client.Do(req)
	return res.Body.String()
}

func (client *Client) TestRouteExpectInvalidRequestFail(t *testing.T, request func() *httptest.ResponseRecorder) {
	// Request
	res := request()

	// Validation
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

type FormData struct {
	Buffer *bytes.Buffer
	writer *multipart.Writer
}

func NewFormData() *FormData {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	return &FormData{
		Buffer: buf,
		writer: writer,
	}
}

func (f *FormData) WriteField(name, value string) error {
	return f.writer.WriteField(name, value)
}

func (f *FormData) WriteFileField(name, path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("filepath cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := f.writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	return err
}

func (f *FormData) Close() error {
	return f.writer.Close()
}

func (f *FormData) ContentType() string {
	return f.writer.FormDataContentType()
}
