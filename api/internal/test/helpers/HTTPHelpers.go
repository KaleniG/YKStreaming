package testhelpers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// * CORS path route should always succeed
func (client *Client) HandleOptionsCORSExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("OPTIONS", "/http/auth/check", nil)
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusNoContent, res.Code)
}
