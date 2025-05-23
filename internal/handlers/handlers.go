// internal/handlers/handlers.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler renders the home page
func IndexHandler(PageTitle string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"PageTitle": "Home",
			//"PageTitle": PageTitle,
			"Content": "Welcome to the home page!",
		})
	}
}

// InitLogsHandler renders the initialization logs
func InitLogsHandler(PageTitle string, initLogs string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"PageTitle": "Init Logs",
			//"PageTitle": PageTitle,
			"logs": initLogs, // Assuming logs is a string you want to display
		})
	}
}
