// internal/routes/routes.go
package routes

import (
	"big_go/internal/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the main routes for the application
func SetupRoutes(r *gin.Engine, PageTitle string) {
	fmt.Println("Setting up routes with page title:", PageTitle)
	r.GET("/", func(c *gin.Context) {
		fmt.Println("Handling request for / %s", PageTitle)
		handlers.IndexHandler(PageTitle)(c)
	})
}

// SetupInitLogsRoute sets up the /init_logs route
func SetupInitLogsRoute(r *gin.Engine, PageTitle, initLogs string) {
	fmt.Println("Setting up /init_logs route with page title:", PageTitle, "and initLogs:", initLogs)

	r.GET("/init_logs", func(c *gin.Context) {
		fmt.Println("Handling request for /init_logs %s", PageTitle)
		handlers.InitLogsHandler(PageTitle, initLogs)(c)
	})
}

// SetupInitLogsRoute sets up the /init_logs route
//func SetupInitLogsRoute(r *gin.Engine, pageTitle, initLogs string) {
//	r.GET("/init_logs", handlers.InitLogsHandler(pageTitle, initLogs))
//}
