package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	formaterror "github.com/buraktabakoglu/GOLANGAPPX/api/pkg/utils"
	//"github.com/gorilla/mux"
)

// Register creates a new user account and sends an activation email.
// @Summary Create a new user account
// @Description creates a new user account by taking user's email, name, password and other details.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User Details"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /api/register [post]
func (server *Server) CreateUser(c *gin.Context) {

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
	errorMessages := user.Validate("login")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		err = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}
	token := auth.RegisterCreateToken(user.Email, user.CreatedAt)
	result := server.DB.Exec("INSERT INTO activation_links(user_id, token, is_used, created_at) VALUES($1, $2, $3, $4)", user.ID, token, false, time.Now())
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save activation link"})
		return
	}

	go func() {
		w := kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{os.Getenv("KAFKA_BROKER")},
			Topic:    "email-topic",
			Balancer: &kafka.LeastBytes{},
		})
		defer w.Close()

		user.Token = token
		userJSON, _ := json.Marshal(user)
		w.WriteMessages(context.Background(),
			kafka.Message{
				Value: userJSON,
			},
		)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": userCreated,
	})
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func (server *Server) GetUsers(c *gin.Context) {
	errList = map[string]string{}

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": users,
	})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [get]
func (server *Server) GetUser(c *gin.Context) {

	errList = map[string]string{}

	userID := c.Param("id")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	user := models.User{}

	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

// UpdateUser godoc
// @Summary Update a User by ID
// @Description Update a User by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param user body models.User true "Update User"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [put]
func (server *Server) UpdateUser(c *gin.Context) {

	errList = map[string]string{}

	userID := c.Param("id")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

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

	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenID != uint32(uid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	user.Prepare()
	errorMessages := user.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})
}

// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [delete]
func (server *Server) DeleteUser(c *gin.Context) {

	var tokenID uint32
	errList = map[string]string{}

	userID := c.Param("id")
	user := models.User{}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenID, err = auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenID != 0 && tokenID != uint32(uid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}



// ActivateUser godoc
// @Summary Activate a User by token
// @Description Activate a User by token
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token path string true "Activation Token"
// @Success 200 
// @Failure 400 
// @Failure 400 
// @Failure 500 
// @Failure 500 
// @Router /api/activate/{token} [get]
func (server *Server) ActivateUser(c *gin.Context) {
	token := c.Param("token")
	

	dbuser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbuser, password, dbname))
	if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
	return
	}
	defer db.Close()


	// check if the activation link already used
	var is_used bool
	var userID int
	err = db.QueryRow("SELECT is_used, user_id FROM activation_links WHERE token = $1", token).Scan(&is_used, &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activation token"})
		return
	}
	if is_used {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activation link already used"})
		return
	}
	// Update the user's activation status in the database
	_, err = db.Exec("UPDATE users SET is_active = true WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate the user"})
		return
	}
	_, err = db.Exec("UPDATE activation_links SET is_used = true WHERE token = $1", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activation link status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User activated successfully"})
}
