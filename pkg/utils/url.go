package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.FullPath())
}