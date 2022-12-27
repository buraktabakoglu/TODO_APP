package test

import (
	"bytes"
	"encoding/json"
	"github.com/buraktabakoglu/GOLANGAPPX/api/internal/mail"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	sendMailFunc func(ToUser string, FromAdmin string, Token string, Sendgridkey string, AppEnv string) (*mail.EmailResponse, error)
)
type sendMailMock struct {}

func (sm *sendMailMock) SendResetPassword(ToUser string, FromAdmin string, Token string, Sendgridkey string, AppEnv string) (*mail.EmailResponse, error) {
	return sendMailFunc(ToUser, FromAdmin, Token, Sendgridkey, AppEnv)
}

func TestForgotPasswordSuccess(t *testing.T) {

	

	gin.SetMode(gin.TestMode)

	err := refreshUserAndResetPasswordTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	
	mail.SendMail = &sendMailMock{} 

	
	sendMailFunc = func(ToUser string, FromAdmin string, Token string, Sendgridkey string, AppEnv string) (*mail.EmailResponse, error) {
		return &mail.EmailResponse{
			Status:   http.StatusOK,
			RespBody: "Success, Please click on the link provided in your email",
		}, nil
	}
		inputJSON :=  `{"email": "naber17@gmail.com"}` //the seeded user
		r := gin.Default()
		r.POST("/password/forgot", server.ForgotPassword)
		req, err := http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBufferString(inputJSON))
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
	
		status := responseInterface["status"]

		assert.Equal(t, rr.Code, int(status.(float64))) 
		
}


func TestForgotPasswordFailures(t *testing.T) {

	

	gin.SetMode(gin.TestMode)

	err := refreshUserAndResetPasswordTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	
	mail.SendMail = &sendMailMock{} 

	samples := []struct {
		id         string
		inputJSON  string
		statusCode int
	}{
		{
			
			inputJSON:  `{"email": "gasgsaexample.com"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"email": "rburakn@example.com"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"email": ""}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"email": 123}`,
			statusCode: 422,
		},
	}
	for _, v := range samples {
		r := gin.Default()
		r.POST("/password/forgot", server.ForgotPassword)
		req, err := http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBufferString(v.inputJSON))
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
		if v.statusCode == 422 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["invalid_email"] != nil {
				assert.Equal(t, responseMap["invalid_email"], "invalid Email")
			}
			if responseMap["no_email"] != nil {
				assert.Equal(t, responseMap["no_email"], "Sorry, we do not recognize this email")
			}
			if responseMap["required_email"] != nil {
				assert.Equal(t, responseMap["required_email"], "required Email")
			}
			if responseMap["unmarshal_error"] != nil {
				assert.Equal(t, responseMap["unmarshal_error"], "Cannot unmarshal body")
			}
		}
	}
}

func TestResetPassword(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndResetPasswordTable()
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedResetPassword()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			
			inputJSON:  `{"token": ""}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "23423498398rwnef9sd8fjsdf"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "awesometoken", "new_password": "pass", "retype_password":"pass"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "awesometoken", "new_password": "", "retype_password":"password"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":""}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":"newpassword"}`,
			statusCode: 422,
		},
		{
			
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":"password"}`,
			statusCode: 200,
		},
	}
	for _, v := range samples {
		r := gin.Default()
		r.POST("/password/reset", server.ResetPassword)
		req, err := http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBufferString(v.inputJSON))
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

		if v.statusCode == 200 {
			responseMap := responseInterface["response"]
			assert.Equal(t, responseMap,"Success")
		}
		if v.statusCode == 422 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_token"] != nil {
				assert.Equal(t, responseMap["Invalid_token"], "Invalid link. Try requesting again")
			}
			if responseMap["No_email"] != nil {
				assert.Equal(t, responseMap["No_email"], "Sorry, we do not recognize this email")
			}
			if responseMap["Invalid_Passwords"] != nil {
				assert.Equal(t, responseMap["Invalid_Passwords"], "Password should be atleast 6 characters")
			}
			if responseMap["Empty_passwords"] != nil {
				assert.Equal(t, responseMap["Empty_passwords"], "Please ensure both field are entered")
			}
			if responseMap["Password_unequal"] != nil {
				assert.Equal(t, responseMap["Password_unequal"], "Passwords provided do not match")
			}
		}
	}
}