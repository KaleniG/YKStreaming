package handlers_test

import (
	"testing"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"
)

func TestOptionsCORSHandler(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupAuthRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		client.HandleOptionsCORSExpectSuccess(t)
	}
}
