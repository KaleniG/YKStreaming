package testhelpers

import (
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// * Simulates an actual stream starting and going trough all the events
func SimulateRealtimeStreamPublishing(t *testing.T, streamKey string) {
	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-f", "lavfi", "-i", "testsrc=size=640x360:rate=25",
		"-f", "lavfi", "-i", "sine",
		"-c:v", "libx264",
		"-tune", "zerolatency",
		"-f", "flv",
		"rtmp://localhost/live/"+streamKey,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Start()
}

// * Call the RTMP OnStreamPublish event and make checks expecting success
// * Success happens only if the stream key of the live exists in the database
func (client *Client) SimulateOnStreamPublishExpectSuccess(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)
}

// * Call the RTMP OnStreamPublish event and make checks expecting fail
// * Fail happens only if the stream key of the live does not exists in the database
func (client *Client) SimulateOnStreamPublishUnexistentExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusNotFound, res.Code)
}

// * Call the RTMP OnStreamPublish duplicate (the stream has already started) event and make checks expecting fail
func (client *Client) SimulateOnStreamPublishDuplicateOrStreamEndedExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusForbidden, res.Code)
}

// * Call the RTMP OnStreamPublishDone event and make checks expecting success
func (client *Client) SimulateOnStreamPublishDoneExpectSuccess(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish-done", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)
}

// * Call the RTMP OnStreamPublishDone event and make checks expecting fail
// * Fail happens only if the stream key of the live does not exists in the database
func (client *Client) SimulateOnStreamPublishDoneUnexistentExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish-done", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusNotFound, res.Code)
}

// * Call the RTMP OnStreamPublishDone duplicate (the stream has already ended) event and make checks expecting fail
func (client *Client) SimulateOnStreamPublishDoneDuplicateOrNotStartedExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-publish-done", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusForbidden, res.Code)
}

// * Call the RTMP OnStreamUpdate event and make checks expecting success
func (client *Client) SimulateOnStreamUpdateExpectSuccess(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-update", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)
}

// * Call the RTMP OnStreamPublishDone event and make checks expecting fail
// * Fail happens only if the stream key of the live does not exists in the database
func (client *Client) SimulateOnStreamUpdateUnexistentExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-update", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusNotFound, res.Code)
}

// * Call the RTMP OnStreamPublishDone duplicate (the stream has already ended) event and make checks expecting fail
func (client *Client) SimulateOnStreamUpdateStreamNotStartedOrEndedExpectFail(t *testing.T, streamKey string) {
	// Setup & request
	formData := NewFormData()
	formData.WriteField("name", streamKey)
	formData.Close()
	req, _ := http.NewRequest("POST", "/rtmp/on-update", formData.Buffer)
	req.Header.Set("Content-Type", formData.ContentType())
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusForbidden, res.Code)
}
