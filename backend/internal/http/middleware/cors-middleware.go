package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders: []string{"Content-Length"},

		AllowCredentials: true,
		AllowOrigins: []string{
			"https://scrumno-api.ru",
			"http://localhost:3001",
		},
	}))
}
