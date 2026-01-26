package handlers_test

import (
	"testing"
	"ykstreaming_api/internal/helpers"
	testhelpers "ykstreaming_api/internal/test/helpers"
	testsetups "ykstreaming_api/internal/test/setups"

	"github.com/stretchr/testify/assert"
)

// * When not logged in or signed in there shall not be any user in session or cookies
func TestCheck(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupAuthRouter()
	defer dbStore.Close()
	client := testhelpers.NewClient(router)

	// Validation
	{
		client.AuthCheckExpectFail(t)
		assert.Empty(t, client.GetCookieKey("remember_token"))
		assert.Empty(t, client.GetSessionKey("user_id"))
	}
}

// * After the signup there shall be a user_id in the session
func TestSignup(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupAuthRouter()
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

		// Second signup, duplicate, should be invalid
		client.SigninDuplicateExpectFail(t, name, email, password)

		// Cookies check
		assert.Empty(t, client.GetCookieKey("remember_token"))
	}
}

// * It starts with a signup and after that two invalid logins,
// * one for each parameter and then two logins one with remember me and one without
func TestLogin(t *testing.T) {
	// Setup
	router, dbStore := testsetups.SetupAuthRouter()
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

		// * Here start the invalid params login cases
		// Case 1: invalid password
		client.LoginInvalidPasswordExpectFail(t, email, password)

		// Auth & cookie check
		client.AuthCheckExpectFail(t)
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Case 2: invalid email
		client.LoginInvalidEmailExpectFail(t, email, password)

		// Auth & cookie check
		client.AuthCheckExpectFail(t)
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// * Here start the valid params login cases with and without remember me option
		// Case 1: Without remember me
		client.LoginExpectSuccess(t, email, password)

		// Cookie check
		assert.Empty(t, client.GetCookieKey("remember_token"))

		// Case 2: With remember me
		client.LoginRememberMeExpectSuccess(t, email, password)

		// Now a logout again to check logout after login
		client.LogoutExpectSuccess(t)
	}
}
