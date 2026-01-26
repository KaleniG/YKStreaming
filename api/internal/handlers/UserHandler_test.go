package handlers_test

import (
	"testing"
	"ykstreaming_api/internal/helpers"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"

	"github.com/stretchr/testify/assert"
)

// * It starts with a signup and after that a logout
func TestLogout(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"
	name := random + "_name"
	password := random + "_password"

	{
		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Now a logout
		client.LogoutExpectSuccess(t)
	}
}

func TestGetUserStreams(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"
	name := random + "_name"
	password := random + "_password"
	streamName := random + "_stream"
	thumbnailPath := "./../../assets/thumbnail.png"

	{
		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Get all streams but none were created so an empty list
		client.GetUserStreamsEmptyExpectSuccess(t)

		// * Add stream cases
		// Case 1: Just add 4 streams and check for their existance
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		streams := client.GetUserStreamsExpectSuccess(t)
		assert.Len(t, streams, 4)

		// Case 2: Start the 4 streams and then stop them while checking that data is coherent
		// Start stream simulation
		for _, key := range streamKeyList {
			client.SimulateOnStreamPublishExpectSuccess(t, key)

			streams = client.GetUserStreamsExpectSuccess(t)
			assert.Len(t, streams, 4)

			for _, stream := range streams {
				if stream.Key == key {
					assert.False(t, stream.StartedAt.IsZero())
					assert.True(t, stream.IsActive)
				}
			}
		}

		// TODO: Simulate viewership

		streams = client.GetUserStreamsExpectSuccess(t)
		assert.Len(t, streams, 4)

		// Start stream stopping
		for _, key := range streamKeyList {
			client.StopStreamExpectSuccess(t, key)

			streams := client.GetUserStreamsExpectSuccess(t)
			assert.Len(t, streams, 4)

			for _, stream := range streams {
				if stream.Key == key {
					assert.Equal(t, int64(0), stream.LiveViewers)
					assert.False(t, stream.EndedAt.IsZero())
					assert.False(t, stream.IsActive)
				}
			}
		}

		streams = client.GetUserStreamsExpectSuccess(t)
		assert.Len(t, streams, 4)

		// Case 3: Remove all 4 streams and then check for their existence
		for _, key := range streamKeyList {
			client.RemoveStreamExpectSuccess(t, key)
		}

		client.GetUserStreamsEmptyExpectSuccess(t)
	}
}

func TestAddStream(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		random, _ := helpers.GenerateRandomToken(5)
		email := random + "@gmail.com"
		name := random + "_name"
		password := random + "_password"
		streamName := random + "_stream"
		thumbnailPath := "./../../assets/thumbnail.png"

		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Add stream cases
		client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath)
		client.AddStreamExpectSuccess(t, streamName, false, nil)
		client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath)
		client.AddStreamExpectSuccess(t, streamName, true, nil)
	}
}

func TestRemoveStream(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"
	name := random + "_name"
	password := random + "_password"
	streamName := random + "_stream"
	thumbnailPath := "./../../assets/thumbnail.png"

	{
		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Check removing inexistent stream
		client.RemoveUnexistentStreamExpectFail(t, random)

		// Add a streams
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		// Remove the streams
		for _, streamKey := range streamKeyList {
			client.RemoveStreamExpectSuccess(t, streamKey)
		}
	}
}

func TestStopStream(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"
	name := random + "_name"
	password := random + "_password"
	streamName := random + "_stream"
	thumbnailPath := "./../../assets/thumbnail.png"

	{
		// First signup, should be valid
		client.SignupExpectSuccess(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Check removing inexistent stream
		client.StopUnexistentStreamExpectFail(t, random)

		// Add a streams
		streamKeyList := []string{
			client.AddStreamExpectSuccess(t, streamName, false, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, false, nil),
			client.AddStreamExpectSuccess(t, streamName, true, &thumbnailPath),
			client.AddStreamExpectSuccess(t, streamName, true, nil),
		}

		// Stop the streams
		for _, streamKey := range streamKeyList {
			client.StopStreamExpectSuccess(t, streamKey)
		}

		// Remove the streams
		for _, streamKey := range streamKeyList {
			client.RemoveStreamExpectSuccess(t, streamKey)
		}
	}
}
