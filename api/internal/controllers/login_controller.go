package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	formaterror "github.com/buraktabakoglu/GOLANGAPPX/api/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login
// @Description Logs in a user and returns a JWT
// @Produce json
// @Param body body models.User true "User email and password"
// @Success 200 string token
// @Router /api/users/login [post]
// @Security ApiKeyAuth
func (server *Server) Login(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":      http.StatusUnprocessableEntity,
			"first error": "Unable to get request",
		})
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Cannot unmarshal body",
		})
		return
	}

	user.Prepare()
	errList = user.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  formattedError,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": token,
		})
	}

}

// SignIn godoc
// @Summary Sign in user and get token
// @Description Get token for authenticated user
// @Produce json
// @Param email query string true "User email"
// @Param password query string true "User password"
// @Success 200 string token
// @Failure 400 string Error
// @Failure 401 string Error
// @Router /api/signin [post]
func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(user.ID))
}

// Logout godoc
// @Summary Logout a user
// @Description Revokes the user's access token
// @Produce text/plain
// @Success 200 {string} string "Old cookie deleted and token blacklisted. Logged out!"
// @Router /api/users/logout [get]
// @Security ApiKeyAuth
func (server *Server) Logout(c *gin.Context) {

	tokenString := auth.ExtractToken(c.Request)

	if err := auth.TokenValid(c.Request); err != nil {

		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tokenHash := auth.TokenHash(tokenString)

	result := server.DB.Exec("INSERT INTO blacklisted_tokens (token) VALUES (?)", tokenHash)
	if err := result.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", "", -1, "", "", false, true)
}
