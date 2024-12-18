package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func URLCleanerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		originalPath := c.Request.URL.Path

		if len(originalPath) > 1000 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		cleanedPath := ""
		for i := 0; i < len(originalPath); i++ {
			if i > 0 && originalPath[i] == '/' && cleanedPath[len(cleanedPath)-1] == '/' {
				continue
			}
			cleanedPath += string(originalPath[i])
		}

		if len(cleanedPath) > 1 && cleanedPath[len(cleanedPath)-1] == '/' {
			cleanedPath = cleanedPath[:len(cleanedPath)-1]
		}

		if cleanedPath != originalPath {
			c.Redirect(http.StatusMovedPermanently, cleanedPath)
			c.Abort()
			return
		}

		c.Next()
	}
}
