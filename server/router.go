package server

import (
	"web-service-gin/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// Initialize default Engine instance
	router := gin.Default()

	// Initialize controller instances
	document := controllers.NewDocumentController()

	// Set up API paths
	router.GET("/document/:oid", document.GetDocument)

	return router
}
