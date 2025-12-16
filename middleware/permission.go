package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hasPermission(list []string, required string) bool {
	for _, p := range list {
		if p == required {
			return true
		}
	}
	return false
}

func Permission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		permissionsRaw, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "permissions not found in token",
			})
			c.Abort()
			return
		}

		permissions, ok := permissionsRaw.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "invalid permission format",
			})
			c.Abort()
			return
		}

		if !hasPermission(permissions, requiredPermission) {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "you do not have permission: " + requiredPermission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
