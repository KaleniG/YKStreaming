package handlers_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"ykstreaming_api/internal/helpers"
	"ykstreaming_api/internal/test"

	"github.com/go-playground/assert/v2"
)

// When not logged in or signed in there shall not be any user in session or cookies
func TestCheck(t *testing.T) {
	router, dbStore := test.SetupAuthRouter()
	defer dbStore.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/http/auth/check", nil)
	router.ServeHTTP(w, req)

	// Status code check
	assert.Equal(t, http.StatusOK, w.Code)

	// Body content check
	var res struct {
		User bool `json:"user"`
	}
	err := json.Unmarshal([]byte(w.Body.String()), &res)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, res.User)

	// Cookies content check ("remember_token" should not be set)
	resCookies := w.Result().Cookies()
	for _, c := range resCookies {
		if c.Name == "remember_token" {
			assert.Equal(t, "", string(c.Value))
		}
	}

	// Session content check
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/http/session/user_id", nil)
	for _, c := range resCookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, "", w.Body.String())
}

// After the signup there shall be a user_id in the session
func TestSignup(t *testing.T) {
	router, dbStore := test.SetupAuthRouter()
	defer dbStore.Close()

	client := test.NewClient(router)

	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"

	body := `{
		"name": "` + random + `",
		"email": "` + email + `",
		"password": "1234"
	}`

	// --- signup ---
	req, _ := http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res := client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- auth check after signup ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	var authRes struct {
		User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &authRes)

	assert.Equal(t, email, authRes.User.Email)

	// --- duplicate signup ---
	req, _ = http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res = client.Do(req)
	assert.Equal(t, http.StatusConflict, res.Code)
	var signupErr struct {
		Param string `json:"param"`
		Error string `json:"error"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &signupErr)

	assert.Equal(t, "email", signupErr.Param)

	// --- auth still valid---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	var authRes2 struct {
		User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
	err := json.Unmarshal(res.Body.Bytes(), &authRes2)
	if err != nil {
		log.Print(err)
	}

	assert.Equal(t, email, authRes2.User.Email)

	// cleanup
	ctx := context.Background()
	dbStore.Queries.RemoveUserByEmail(ctx, email)
}

// It starts with a signup and after that a logout
func TestLogout(t *testing.T) {
	router, dbStore := test.SetupAuthRouter()
	defer dbStore.Close()

	client := test.NewClient(router)

	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"

	body := `{
		"name": "` + random + `",
		"email": "` + email + `",
		"password": "1234"
	}`

	// --- signup ---
	req, _ := http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res := client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- logout ---
	req, _ = http.NewRequest("POST", "/http/user/logout", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- auth still invalid ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	var checkRes struct {
		User bool `json:"user"`
	}
	_ = json.Unmarshal([]byte(res.Body.String()), &checkRes)
	assert.Equal(t, false, checkRes.User)

	// cleanup
	ctx := context.Background()
	dbStore.Queries.RemoveUserByEmail(ctx, email)
}

// It starts with a signup and after that two invalid logins,
// one for each parameter and then two logins one with remember me and one without
func TestLogin(t *testing.T) {
	router, dbStore := test.SetupAuthRouter()
	defer dbStore.Close()

	client := test.NewClient(router)

	random, _ := helpers.GenerateRandomToken(5)
	email := random + "@gmail.com"

	body := `{
		"name": "` + random + `",
		"email": "` + email + `",
		"password": "1234"
	}`

	// --- signup ---
	req, _ := http.NewRequest("POST", "/http/auth/signup", strings.NewReader(body))
	res := client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- logout ---
	req, _ = http.NewRequest("POST", "/http/user/logout", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- login (invalid password case) ---
	body = `{
		"email": "` + email + `",
		"password": "12345asds",
		"remember_me": false
	}`

	req, _ = http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res = client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	var loginResError1 struct {
		Param string `json:"param"`
		Error string `json:"error"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &loginResError1)
	assert.Equal(t, "password", loginResError1.Param)

	// --- login (invalid email case) ---
	body = `{
		"email": "dsajshhd` + email + `",
		"password": "1234",
		"remember_me": false
	}`

	req, _ = http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res = client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	var loginResError2 struct {
		Param string `json:"param"`
		Error string `json:"error"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &loginResError2)
	assert.Equal(t, "email", loginResError2.Param)

	// --- auth still invalid ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	var checkRes struct {
		User bool `json:"user"`
	}
	_ = json.Unmarshal([]byte(res.Body.String()), &checkRes)
	assert.Equal(t, false, checkRes.User)

	// --- login (success case but without remembering the user) ---
	body = `{
		"email": "` + email + `",
		"password": "1234",
		"remember_me": false
	}`

	req, _ = http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- auth valid ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)

	var authRes2 struct {
		User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &authRes2)

	assert.Equal(t, email, authRes2.User.Email)

	// --- no remember_token in cookies ---
	resCookies := res.Result().Cookies()
	for _, c := range resCookies {
		if c.Name == "remember_token" {
			assert.Equal(t, "", string(c.Value))
		}
	}

	// --- login (success case but remembering the user) ---
	body = `{
		"email": "` + email + `",
		"password": "1234",
		"remember_me": true
	}`

	req, _ = http.NewRequest("POST", "/http/auth/login", strings.NewReader(body))
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- auth valid ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)

	var authRes3 struct {
		User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &authRes3)

	assert.Equal(t, email, authRes3.User.Email)

	// --- remember_token in cookies ---
	resCookies = res.Result().Cookies()
	for _, c := range resCookies {
		if c.Name == "remember_token" {
			assert.NotEqual(t, "", string(c.Value))
		}
	}

	// --- logout ---
	req, _ = http.NewRequest("POST", "/http/user/logout", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	// --- auth still invalid ---
	req, _ = http.NewRequest("POST", "/http/auth/check", nil)
	res = client.Do(req)
	assert.Equal(t, http.StatusOK, res.Code)

	var checkRes3 struct {
		User bool `json:"user"`
	}
	_ = json.Unmarshal([]byte(res.Body.String()), &checkRes3)
	assert.Equal(t, false, checkRes3.User)

	// cleanup
	ctx := context.Background()
	dbStore.Queries.RemoveUserByEmail(ctx, email)
}
