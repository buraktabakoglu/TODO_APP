package middlewares

import (
	"net/http"
	"os"
	"strconv"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/gin-gonic/gin"
)

func CombinedAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")

		url := os.Getenv("AUTHORIZE_URL")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Not authorized"})
			c.Abort()
			return
		}

		

		c.Next()
	}
}


func CheckUserOwnership(c *gin.Context) {
	userID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Error in extracting user ID from token"})
		c.Abort()
		return
	}
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		c.Abort()
		return
	}

	if userID != uint32(contentID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This content does not belong to the"})
		c.Abort()
		return
	}
}
