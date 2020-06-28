package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}

	samples := []struct {
		nickname     string
		password     string
		errorMessage string
	}{
		{
			nickname:     user.Nickname,
			password:     "password", //Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			nickname:     user.Nickname,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			nickname:     "Wrong nickname",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {

		token, err := server.SignIn(v.nickname, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {

	refreshUserTable()

	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"nickname": "Steven", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"nickname": "Marcelo", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			inputJSON:    `{"nickname": "Frank", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Details",
		},
		{
			inputJSON:    `{"nickname": "Nathan", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"nickname": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Nickname",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
