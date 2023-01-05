package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	

		tokenString := auth.ExtractToken(c.Request)
		tokenHash := auth.TokenHash(tokenString)

		redisConn := auth.GetRedisConnection()
		isBlacklisted, err := redisConn.SIsMember("blacklisted_tokens", tokenHash).Result()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return 
		}

		if isBlacklisted {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token geçersizdir"))
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
