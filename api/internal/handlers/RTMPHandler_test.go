package handlers_test

import (
	"testing"
	"time"
	"ykstreaming_api/internal/helpers"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"

	"github.com/stretchr/testify/assert"
)

func TestOnPublish(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupRTMPRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		random, _ := helpers.GenerateRandomToken(24)
		client.SimulateOnStreamPublishUnexistentExpectFail(t, random)

		random, _ = helpers.GenerateRandomToken(5)
		email := random + "@gmail.com"
		name := random + "_name"
		password := random + "_password"
		streamName := random + "_stream"
		thumbnailPath := "./../../assets/thumbnail.png"

		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Add streams
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		// Start each stream
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishExpectSuccess(t, key)
		}

		// Try starting again
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishDuplicateOrStreamEndedExpectFail(t, key)
		}

		// Stop the streams
		for _, key := range streamKeyList {
			client.StopStreamExpectSuccess(t, key)
		}

		// Try starting again
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishDuplicateOrStreamEndedExpectFail(t, key)
		}
	}
}

func TestOnPublishDone(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupRTMPRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		// Try to trigger when the stream doesnt even exist
		random, _ := helpers.GenerateRandomToken(24)
		client.SimulateOnStreamPublishDoneUnexistentExpectFail(t, random)

		random, _ = helpers.GenerateRandomToken(5)
		email := random + "@gmail.com"
		name := random + "_name"
		password := random + "_password"
		streamName := random + "_stream"
		thumbnailPath := "./../../assets/thumbnail.png"

		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Add streams
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		// Try to trigger when the streams exists but hasnt started yet
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishDoneDuplicateOrNotStartedExpectFail(t, key)
		}

		// Start each stream
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishExpectSuccess(t, key)
		}

		// Trigger when streams exist and are started
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishDoneExpectSuccess(t, key)
		}

		// Try to trigger again (duplicate)
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishDoneDuplicateOrNotStartedExpectFail(t, key)
		}
	}
}

func TestOnStreamUpdate(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupRTMPServerRouter()
	defer dbStore.Close()
	//defer server.Process.Kill()
	client := testhelpers.NewClient(router)

	// Validate
	{
		random, _ := helpers.GenerateRandomToken(5)
		email := random + "@gmail.com"
		name := random + "_name"
		password := random + "_password"
		streamName := random + "_stream"

		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Try before stream exists
		random, _ = helpers.GenerateRandomToken(24)
		client.SimulateOnStreamUpdateUnexistentExpectFail(t, random)

		// Add streams
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		// Try before stream started
		for _, key := range streamKeyList {
			client.SimulateOnStreamUpdateStreamNotStartedOrEndedExpectFail(t, key)
		}

		// Simulate real RTMP stream
		for _, key := range streamKeyList {
			testhelpers.SimulateRealtimeStreamPublishing(t, key)
		}

		time.Sleep(4 * time.Second)

		// Try after stream started
		for _, key := range streamKeyList {
			client.SimulateOnStreamUpdateExpectSuccess(t, key)
		}

		// Try after more than a minute
		//time.Sleep(1*time.Minute + 5*time.Second)
		for _, key := range streamKeyList {
			client.SimulateOnStreamUpdateExpectSuccess(t, key)
		}

		// Stop the stream
		for _, key := range streamKeyList {
			client.StopStreamExpectSuccess(t, key)
		}

		// Try after stream ended
		for _, key := range streamKeyList {
			client.SimulateOnStreamUpdateStreamNotStartedOrEndedExpectFail(t, key)
		}
	}
}

func TestOnStreamRecordDone(t *testing.T) {

}
