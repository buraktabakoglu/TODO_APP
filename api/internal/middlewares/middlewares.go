package middlewares

import (
	"net/http"
	"strconv"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/gin-gonic/gin"
)




func TokenAuthMiddleware() gin.HandlerFunc {
	errList := make(map[string]string)
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			
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
