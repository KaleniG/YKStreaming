package testhelpers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// * Call to get all public streams when none were created and make checks expecting success
// * In this case success is an empty list, which is valid comparing to AuthCheckExpectFail,
// * which is needed to check wether a this is false or true
// * It is not expected that there is a specific type of authed user so no internal auth checks
func (client *Client) GetStreamsEmptyExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/get-streams", nil)
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)

	// Body content check
	var response struct {
		Streams bool `json:"streams"`
	}
	err := json.Unmarshal([]byte(res.Body.String()), &response)
	assert.Nil(t, err)
	assert.False(t, response.Streams)
}

// * Call to get all public streams and make checks expecting success
// * It is not expected that there is a specific type of authed user so no internal auth checks
func (client *Client) GetStreamsExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/get-streams", nil)
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)

	// Body content check
	var response struct {
		Streams []struct {
			StreamerName       string `json:"streamer_name"`
			Key                string `json:"key"`
			Name               string `json:"name"`
			HasCustomThumbnail bool   `json:"has_custom_thumbnail"`
			IsLive             bool   `json:"is_live"`
			IsVod              bool   `json:"is_vod"`
			LiveViewers        int    `json:"live_viewers"`
		} `json:"streams"`
	}
	err := json.Unmarshal([]byte(res.Body.String()), &response)
	assert.Nil(t, err)
	assert.Greater(t, len(response.Streams), 0)
	for _, stream := range response.Streams {
		assert.NotEmpty(t, stream.Name)
		assert.NotEmpty(t, stream.Key)
		assert.True(t, stream.IsLive || stream.IsVod)
		assert.GreaterOrEqual(t, stream.LiveViewers, 0)
	}
}
