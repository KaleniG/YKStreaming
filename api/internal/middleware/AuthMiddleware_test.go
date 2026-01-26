package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"ykstreaming_api/internal/helpers"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"

	"github.com/stretchr/testify/assert"
)

// * Runs UserRoute routes with and without user auth
// * In the initial TestAuthMiddlewareExpectFail cases no request body is used
// * since the requests should be refuted before the body content check even happens
func TestAuthMiddleware(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupUserRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		// * No auth cases
		{
			// Logout
			client.TestAuthMiddlewareExpectFail(t, func() *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/http/user/logout", nil)
				return client.Do(req)
			})

			// Streams
			client.TestAuthMiddlewareExpectFail(t, func() *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/http/user/streams/", nil)
				return client.Do(req)
			})

			// Add Stream
			client.TestAuthMiddlewareExpectFail(t, func() *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/http/user/streams/add", nil)
				return client.Do(req)
			})

			// Remove Stream
			client.TestAuthMiddlewareExpectFail(t, func() *httptest.ResponseRecorder {
				random, _ := helpers.GenerateRandomToken(5)
				req, _ := http.NewRequest("POST", "/http/user/streams/remove/"+random, nil)
				return client.Do(req)
			})

			// Stop Stream
			client.TestAuthMiddlewareExpectFail(t, func() *httptest.ResponseRecorder {
				random, _ := helpers.GenerateRandomToken(5)
				req, _ := http.NewRequest("POST", "/http/user/streams/stop/"+random, nil)
				return client.Do(req)
			})
		}

		// * Auth cases
		{
			random, _ := helpers.GenerateRandomToken(5)
			email := random + "@gmail.com"
			name := random + "_name"
			password := random + "_password"

			// Sign up
			client.SignupExpectSuccess(t, name, email, password)

			// Cookies check
			assert.Empty(t, client.GetCookieKey("remember_token"))

			// Streams
			client.GetUserStreamsEmptyExpectSuccess(t)

			// Add Stream
			client.TestRouteExpectInvalidRequestFail(t, func() *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/http/user/streams/add", nil)
				return client.Do(req)
			})

			// Remove Stream
			client.RemoveUnexistentStreamExpectFail(t, random)

			// Stop Stream
			client.StopUnexistentStreamExpectFail(t, random)

			// Logout
			client.LogoutExpectSuccess(t)
		}
	}
}
