package handlers_test

import (
	"testing"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"

	"github.com/stretchr/testify/assert"
)

// * Getting all the available streams
func TestGetStreams(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupDefaultRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		// First check that there are no streams
		client.GetStreamsEmptyExpectSuccess(t)

		// Check cookies & session
		client.AuthCheckExpectFail(t)
		assert.NotEmpty(t, client.GetSessionKey("guest_token"))
	}
}

func TestGetStream(t *testing.T) {

}

func TestViewStream(t *testing.T) {

}

func TestUnviewStream(t *testing.T) {

}
