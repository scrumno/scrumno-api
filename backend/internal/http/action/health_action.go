package action

import "github.com/gin-gonic/gin"

func HealthAction(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong (from Go!)",
	})
}
