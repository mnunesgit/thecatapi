package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"../../api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateCat(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Nickname, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		breed        string
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"breed":"The breed" }`,
			statusCode:   201,
			tokenGiven:   tokenString,
			breed:        "The breed",
			errorMessage: "",
		},
		{
			inputJSON:    `{"breed":"The breed" }`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Breed Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"breed":"When no token is passed"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"breed":"When incorrect token is passed"}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"breed": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Breed",
		},
		{
			inputJSON:    `{"breed": "This is a breed"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"breed": "This is an awesome breed"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"breed": "This is an awesome breed"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/breeds", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateCat)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["breed"], v.breed)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetCats(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndCats()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/breeds", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetCats)
	handler.ServeHTTP(rr, req)

	var cats []models.Cat
	err = json.Unmarshal([]byte(rr.Body.String()), &cats)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(cats), 2)
}
func TestGetCatByID(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatal(err)
	}
	cat, err := seedOneUserAndOneCat()
	if err != nil {
		log.Fatal(err)
	}
	catSample := []struct {
		id           string
		statusCode   int
		breed        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(cat.ID)),
			statusCode: 200,
			breed:      cat.Breed,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range catSample {

		req, err := http.NewRequest("GET", "/breeds", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetCat)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, cat.Breed, responseMap["breed"])
		}
	}
}
