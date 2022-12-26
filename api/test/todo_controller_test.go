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

func TestCreateTodo(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface 
	tokenString := fmt.Sprintf("Bearer %v", token)

	
	samples := []struct {
		inputJSON  string
		statusCode int
		status      string
		description    string
		tokenGiven string
	}{
		{
			inputJSON:  `{"status":"The status", "description": "the description"}`,
			statusCode: 201,
			tokenGiven: tokenString,
			status:      "The status",
			description:    "the description",
		},
		{
			
			inputJSON:  `{"status":"The status", "description": "the description"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			
			inputJSON:  `{"status":"When no token is passed", "description": "the description"}`,
			statusCode: 401,
			tokenGiven: "",
		},
		{
			
			inputJSON:  `{"status":"When incorrect token is passed", "description": "the description"}`,
			statusCode: 401,
			tokenGiven: "This is an incorrect token",
		},
		}

	for _, v := range samples {

		r := gin.Default()

		r.POST("/todos", server.CreateTodo)
		req, err := http.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(v.inputJSON))
		req.Header.Set("Authorization", v.tokenGiven)
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
			assert.Equal(t, responseMap["status"], v.status)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["status"] != nil {
				assert.Equal(t, responseMap["status"], "Title Already Taken")
			}
			if responseMap["status"] != nil {
				assert.Equal(t, responseMap["Required_status"], "Required status")
			}
			if responseMap["Required_description"] != nil {
				assert.Equal(t, responseMap["Required_description"], "Required description")
			}
		}
	}
}

func TestGetTodos(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndTodos()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/todos", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/todos", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	todosInterface := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &todosInterface)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	
	theTodos := todosInterface["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theTodos), 2)
}



func TestUpdateTodo(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var TodoUserEmail, TodoUserPassword string
	
	var AuthTodoID uint64

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatal(err)
	}
	users, todos, err := seedUsersAndTodos()
	if err != nil {
		log.Fatal(err)
	}
	
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		TodoUserEmail = user.Email
		TodoUserPassword = "password" 
	}
	
	for _, todo := range todos {
		if todo.ID == 2 {
			continue
		}
		AuthTodoID = todo.ID
	}
	
	tokenInterface, err := server.SignIn(TodoUserEmail, TodoUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface 
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id         string
		updateJSON string
		statusCode int
		status      string
		description    string
		tokenGiven string
	}{
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"The updated todo", "description": "This is the updated description"}`,
			statusCode: 200,
			status:      "",
			description:    "",
			tokenGiven: tokenString,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"This is still another status", "description": "This is the updated description"}`,
			tokenGiven: "",
			statusCode: 401,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"This is still another status", "description": "This is the updated description"}`,
			tokenGiven: "this is an incorrect token",
			statusCode: 401,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"status 2", "description": "This is the updated description"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"", "description": "This is the updated description"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			updateJSON: `{"status":"Awesome status", "description": ""}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range samples {

		r := gin.Default()

		r.PUT("/todos/:id", server.UpdateATodo)
		req, err := http.NewRequest(http.MethodPut, "/todos/"+v.id, bytes.NewBufferString(v.updateJSON))
		req.Header.Set("Authorization", v.tokenGiven)
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

		if v.statusCode == 204 {
			
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["status"], v.status)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 400 || v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Taken_status"] != nil {
				assert.Equal(t, responseMap["Taken_status"], "status Already Taken")
			}
			if responseMap["Required_status"] != nil {
				assert.Equal(t, responseMap["Required_title"], "Required Title")
			}
			if responseMap["Required_description"] != nil {
				assert.Equal(t, responseMap["Required_description"], "Required description")
			}
		}
	}
}

func TestDeleteTodo(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var TodoUserEmail, TodoUserPassword string
	// var AuthID uint32
	var AuthTodoID uint64

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatal(err)
	}
	users, todos, err := seedUsersAndTodos()
	if err != nil {
		log.Fatal(err)
	}
	
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		TodoUserEmail = user.Email
		TodoUserPassword = "password" 
	}
	
	for _, todo := range todos {
		if todo.ID == 1 {
			continue
		}
		AuthTodoID = todo.ID
	}
	
	tokenInterface, err := server.SignIn(TodoUserEmail, TodoUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface 
	tokenString := fmt.Sprintf("Bearer %v", token)

	todoSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			tokenGiven: tokenString,
			statusCode: 200,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			tokenGiven: "",
			statusCode: 401,
		},
		{
			
			id:         strconv.Itoa(int(AuthTodoID)),
			tokenGiven: "This is an incorrect token",
			statusCode: 401,
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range todoSample {
		r := gin.Default()
		r.DELETE("/todos/:id", server.DeleteATodo)
		req, _ := http.NewRequest(http.MethodDelete, "/todos/"+v.id, nil)
		req.Header.Set("Authorization", v.tokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})

		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseInterface["response"], "Todo deleted")
		}

		if v.statusCode == 400 || v.statusCode == 401 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
		}
	}
}