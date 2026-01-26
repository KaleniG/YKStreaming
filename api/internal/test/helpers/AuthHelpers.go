package testhelpers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// * Call to auth check that should fail
func (client *Client) AuthCheckExpectFail(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/auth/check", nil)
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)

	// Body content check
	var response struct {
		User bool `json:"user"`
	}
	err := json.Unmarshal([]byte(res.Body.String()), &response)
	assert.Nil(t, err)
	assert.False(t, response.User)
}

// * Call to auth check that should succeed
func (client *Client) AuthCheckExpectSuccess(t *testing.T) {
	// Setup & request
	req, _ := http.NewRequest("POST", "/http/auth/check", nil)
	res := client.Do(req)

	// Status code check
	assert.Equal(t, http.StatusOK, res.Code)

	// Body content check
	var response struct {
		User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
	err := json.Unmarshal([]byte(res.Body.String()), &response)
	assert.Nil(t, err)
	assert.NotEmpty(t, strings.TrimSpace(response.User.Email))
	assert.NotEmpty(t, strings.TrimSpace(response.User.Name))
}

// * Create intial account and make checks expecting success
func (client *Client) SignupExpectSuccess(t *testing.T, name string, email string, password string) {
	// Setup & request
	body := `{
		"name": "` + name + `",
		"email": "` + email + `",
		"password": "` + password + `"
	}`
	req, _ := http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusOK, res.Code)

		// Whatever the case a valid signup should make the user authenticate
		client.AuthCheckExpectSuccess(t)
	}
}

// * Create duplicate account and make checks expecting failure
// * As an failure case it could happen at any point of user experience so there is a possibility of the user being logged in even though
// * one wouldn't be usually be able to signin before logging out, so the auth check (i.e. AuthCheckExpectXXX) is situation dependent
func (client *Client) SigninDuplicateExpectFail(t *testing.T, name string, email string, password string) {
	// Setup & request
	body := `{
		"name": "` + name + `",
		"email": "` + email + `",
		"password": "` + password + `"
	}`
	req, _ := http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusConflict, res.Code)
		var response struct {
			Param string `json:"param"`
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.Equal(t, "email", strings.TrimSpace(response.Param))
		assert.NotEmpty(t, strings.TrimSpace(response.Error))
	}
}

// * Login with invalid email and make checks expecting failure
// * When the params are invalid there should be no point in checking remember me value since there won't be any login
// * As an failure case it could happen at any point of user experience so there is a possibility of the user being logged in even though
// * one wouldn't be usually be able to login before logging out, so the auth check (i.e. AuthCheckExpectXXX) is situation dependent
func (client *Client) LoginInvalidEmailExpectFail(t *testing.T, email string, password string) {
	// Setup & request
	body := `{
		"email": "invalid_` + email + `",
		"password": "` + password + `",
		"remember_me": false
	}`

	req, _ := http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		var response struct {
			Param string `json:"param"`
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.Equal(t, "email", response.Param)
		assert.NotEmpty(t, response.Error)
	}
}

// * Login with invalid password and make checks expecting failure
// * When the params are invalid there should be no point in checking remember me value since there won't be any login
// * As an failure case it could happen at any point of user experience so there is a possibility of the user being logged in even though
// * one wouldn't be usually be able to login before logging out, so the auth check (i.e. AuthCheckExpectXXX) is situation dependent
func (client *Client) LoginInvalidPasswordExpectFail(t *testing.T, email string, password string) {
	// Setup & request
	body := `{
		"email": "` + email + `",
		"password": "invalid_` + password + `",
		"remember_me": false
	}`
	req, _ := http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		var response struct {
			Param string `json:"param"`
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(res.Body.String()), &response)
		assert.Nil(t, err)
		assert.Equal(t, "password", response.Param)
		assert.NotEmpty(t, response.Error)
	}
}

// * Login with valid credentials and make checks expecting success
// * This one is with remember me set to false, the inexistence of the remember token this way is not guaranteed,
// * it is situation dependent, but the user should be authenticated
func (client *Client) LoginExpectSuccess(t *testing.T, email string, password string) {
	// Setup & request
	body := `{
		"email": "` + email + `",
		"password": "` + password + `",
		"remember_me": false
	}`
	req, _ := http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusOK, res.Code)
		// Whatever the case a valid login should make the user authenticate
		client.AuthCheckExpectSuccess(t)
	}
}

// * Login with valid credentials and make checks expecting success
// * This one is with remember me set to true, the existence of the remember token this way is guaranteed
func (client *Client) LoginRememberMeExpectSuccess(t *testing.T, email string, password string) {
	// Setup & request
	body := `{
		"email": "` + email + `",
		"password": "` + password + `",
		"remember_me": true
	}`
	req, _ := http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res := client.Do(req)

	// Validation
	{
		assert.Equal(t, http.StatusOK, res.Code)

		// Whatever the case a valid login should make the user authenticate and same the cookie remember token
		client.AuthCheckExpectSuccess(t)
		assert.NotEmpty(t, client.GetCookieKey("remember_token"))
	}
}
