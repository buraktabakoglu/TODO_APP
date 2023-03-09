package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Forgot Password
// @Description Send a password reset link to user's email
// @Accept  json
// @Produce json
// @Param body body models.User true "User email"
// @Success 200 
// @Router /api/forgotpassword [post]
func (server *Server) ForgotPassword(c *gin.Context) {

	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user.Prepare()
	errorMessages := user.Validate("forgotpassword")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		errList["No_email"] = "Sorry, we do not recognize this email"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	resetPassword := models.ResetPassword{}
	resetPassword.Prepare()

	token := auth.TokenHash(user.Email)
	resetPassword.Email = user.Email
	resetPassword.Token = token

	resetDetails, err := resetPassword.SaveDatails(server.DB)
	if err != nil {
		errList["incorrect details"] = "incorrect details"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	fmt.Println("THIS OCCURRED HERE")

	resetPass := models.ResetPassword{
		Email: resetDetails.Email,
		Token: resetDetails.Token,
	}
	resetPassJson, err := json.Marshal(resetPass)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}

	// Send message to kafka topic
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    "email-topic",
		Balancer: &kafka.LeastBytes{},
	})
	defer kafkaWriter.Close()

	kafkaMessage := kafka.Message{
		Key: []byte("reset_password"),
		Value: resetPassJson,
	}
	err = kafkaWriter.WriteMessages(context.Background(), kafkaMessage)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Success",
	})
}

// @Summary Reset Password
// @Description Resets user's password with a token
// @Accept  json
// @Produce json
// @Param token path string true "Password reset token"
// @Param body body map[string]string true "New password and retype password"
// @Success 200
// @Router /api/resetpassword/{token} [post]
func (server *Server) ResetPassword(c *gin.Context) {
	errList := map[string]string{}

    token := c.Param("token")

    // Check if the token exists in the reset password table
    resetPassword := models.ResetPassword{}
    err := server.DB.Debug().Model(models.ResetPassword{}).Where("token = ?", token).Take(&resetPassword).Error
    if err != nil {
        errList["Invalid_token"] = "Invalid link. Try requesting again"
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "status": http.StatusUnprocessableEntity,
            "error":  errList,
        })
        return
    }

    // Read the request body
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        errList["Invalid_body"] = "Unable to get request"
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "status": http.StatusUnprocessableEntity,
            "error":  errList,
        })
        return
    }

    requestBody := map[string]string{}
    err = json.Unmarshal(body, &requestBody)
    if err != nil {
        errList["Unmarshal_error"] = "Cannot unmarshal body"
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "status": http.StatusUnprocessableEntity,
            "error":  errList,
        })
        return
    }

	
	if requestBody["new_password"] == "" || requestBody["retype_password"] == "" {
		errList["Empty_passwords"] = "Please ensure both fields are entered"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// Check if the new password and retype password fields have at least 6 characters
	if len(requestBody["new_password"]) < 6 || len(requestBody["retype_password"]) < 6 {
		errList["Invalid_Passwords"] = "Password should be at least 6 characters"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	if requestBody["new_password"] != requestBody["retype_password"] {
		errList["Password_unequal"] = "Passwords provided do not match"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	
	newPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody["new_password"]), bcrypt.DefaultCost)
	if err != nil {
		errList["Error_hashing_password"] = "Error hashing new password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Update the password in the users table
	server.DB.Model(models.User{}).Where("email = ?", resetPassword.Email).Updates(models.User{Password: string(newPassword)})
	// Delete the reset_password token
	_, err = resetPassword.DeleteDatails(server.DB)
	if err != nil {
		errList["Cannot_delete"] = "Cannot delete record, please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Password reset successful",
	})
}
