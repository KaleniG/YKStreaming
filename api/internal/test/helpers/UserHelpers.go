package testhelpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// * Any request practiced inside this function should be bound to auth middleware
func (client *Client) TestAuthMiddlewareExpectFail(t *testing.T, request func() *httptest.ResponseRecorder) {
	// Request
	res := request()

	// Validation
	{
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		// Wathever the case, if the auth middleware fails then there is no authentication no user session and cookies
		client.AuthCheckExpectFail(t)
		assert.Empty(t, client.GetCookieKey("remember_token"))
	}
}

// * Logout of the account and make check expecting success
func (client *Client) LogoutExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/logout", nil)
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusOK, res.Code)

		// Wathever the case logout should invalidate user authentication and remove all session and cookies of the user (except for guest_token)
		client.AuthCheckExpectFail(t)
		assert.Empty(t, client.GetCookieKey("remember_token"))
	}
}

// * Add a stream and make checks expecting success
func (client *Client) AddStreamExpectSuccess(t *testing.T, name string, isVOD bool, thumbnailPath *string) string {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", name)
	formData.WriteField("is_vod", strconv.FormatBool(isVOD))
	if thumbnailPath != nil {
		formData.WriteFileField("thumbnail", *thumbnailPath)
	}
	formData.Close()

	req, _ := http.NewRequest("POST", "/http/user/streams/add", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusOK, res.Code)

		var response struct {
			Key string `json:"key"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.NotEmpty(t, strings.TrimSpace(response.Key))
		assert.Len(t, strings.TrimSpace(response.Key), 24)
		return response.Key
	}
}

// * Call to remove a user stream that does not exist and make checks expecting fail
func (client *Client) RemoveUnexistentStreamExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/remove/"+streamKey, nil)
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusNotFound, res.Code)
		var response struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.NotEmpty(t, response.Error)
	}
}

// * Call to remove a user stream that does not exist and make checks expecting fail
func (client *Client) RemoveStreamExpectSuccess(t *testing.T, streamKey string) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/remove/"+streamKey, nil)
	res := client.Do(req)

	// Validation
	assert.Equal(t, http.StatusOK, res.Code)
}

// * Call to stop a user stream that does not exist and make checks expecting fail
func (client *Client) StopUnexistentStreamExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/remove/"+streamKey, nil)
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusNotFound, res.Code)
		var response struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.NotEmpty(t, response.Error)
	}
}

// * Call to remove a user stream that does not exist and make checks expecting fail
func (client *Client) StopStreamExpectSuccess(t *testing.T, streamKey string) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/stop/"+streamKey, nil)
	res := client.Do(req)

	// Validation
	assert.Equal(t, http.StatusOK, res.Code)
}

// * Call to get all user streams when none were created and make checks expecting success
// * In this case success is an empty list, which is valid comparing to AuthCheckExpectFail,
// * which is needed to check wether a this is false or true
func (client *Client) GetUserStreamsEmptyExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/", nil)
	res := client.Do(req)

	// Validation
	{
		// Status code checks
		assert.Equal(t, http.StatusOK, res.Code)

		// Body content check
		var response struct {
			Streams bool `json:"streams"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.False(t, response.Streams)
	}
}

// * Call to get all user streams when some were created and make checks expecting success
// * Type validation is performed only since all the otehr checks are up to the situation
type UserStreamData struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	IsActive    bool      `json:"is_active"`
	EndedAt     time.Time `json:"ended_at"`
	StartedAt   time.Time `json:"started_at"`
	TotalViews  int64     `json:"total_views"`
	IsVod       bool      `json:"is_vod"`
	LiveViewers int64     `json:"live_viewers"`
}

func (client *Client) GetUserStreamsExpectSuccess(t *testing.T) []UserStreamData {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/user/streams/", nil)
	res := client.Do(req)

	// Validation
	{
		// Status code checks
		assert.Equal(t, http.StatusOK, res.Code)

		// Body content check
		var response struct {
			Streams []UserStreamData `json:"streams"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.Greater(t, len(response.Streams), 0)
		for _, stream := range response.Streams {
			assert.NotEmpty(t, stream.Name)
			assert.NotEmpty(t, stream.Key)
			assert.GreaterOrEqual(t, stream.TotalViews, int64(0))
			assert.GreaterOrEqual(t, stream.LiveViewers, int64(0))
		}

		return response.Streams
	}
}
