package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
		nickname   string
		email      string
	}{
		{
			inputJSON:  `{"nickname":"Pet", "email": "pet@example.com", "password": "password"}`,
			statusCode: 201,
			nickname:   "Pet",
			email:      "pet@example.com",
		},
		{
			inputJSON:  `{"nickname":"Frank", "email": "pet@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"nickname":"Pet", "email": "grand@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"nickname":"Kan", "email": "kanexample.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"nickname": "Pet", "email": "kan@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"nickname": "Kan", "email": "kan@example.com", "password": ""}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/users", server.CreateUser)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Taken_email"] != nil {
				assert.Equal(t, responseMap["Taken_email"], "Email Already Taken")
			}
			if responseMap["Taken_nickname"] != nil {
				assert.Equal(t, responseMap["Taken_nickname"], "nickname Already Taken")
			}
			if responseMap["invalid_email"] != nil {
				assert.Equal(t, responseMap["invalid_email"], "invalid Email")
			}
			if responseMap["Required_nickname"] != nil {
				assert.Equal(t, responseMap["Required_nickname"], "Required nickname")
			}
			if responseMap["required_email"] != nil {
				assert.Equal(t, responseMap["required_email"], "required Email")
			}
			if responseMap["required_password"] != nil {
				assert.Equal(t, responseMap["required_password"], "required Password")
			}
		}
	}
}

func TestGetUsers(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_,err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	usersMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &usersMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	
	theUsers := usersMap["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUsers), 2)
}


func TestDeleteUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface 
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id         string
		tokenGiven string
		statusCode int
	}{
		{
			
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: tokenString,
			statusCode: 200,
		},
		{
			
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: "",
			statusCode: 401,
		},
		{
			
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: "This is an incorrect token",
			statusCode: 401,
		},
		{
			
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		r := gin.Default()
		r.DELETE("/users/:id", server.DeleteUser)
		req, _ := http.NewRequest(http.MethodDelete, "/users/"+v.id,bytes.NewBufferString(v.id))
		req.Header.Set("Authorization", v.tokenGiven)

		fmt.Println("FORM REQUEST: ", req.Header)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal((rr.Body.Bytes()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t,rr.Code,v.statusCode)
		}

		if v.statusCode == 400 || v.statusCode == 401 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
		}
	}
}