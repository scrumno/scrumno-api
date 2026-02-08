package middleware

import "github.com/gin-gonic/gin"

func InitMiddlewares(r *gin.Engine) {
	CorsMiddleware(r)
}
