package main

import (
	"autopost/internal/http/action"
	"autopost/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	middleware.InitMiddlewares(r)

	api := r.Group("/api/v1")
	{
		healthGroup := api.Group("/health")
		{
			healthGroup.GET("/check", action.HealthAction)
		}

		postGroup := api.Group("/posts")
		{
			postGroup.GET("/", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "posts (from Go!)",
				})
			})
		}
	}

	err := r.Run()
	if err != nil {
		return
	}
}
