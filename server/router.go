package server

import (
	"web-service-gin/controllers"
	"web-service-gin/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Initialize default Engine instance
	router := gin.Default()

	// Initialize controller instances
	document := &controllers.DocumentController{DB: database.GetDB()}

	// Set up API paths
	router.GET("/data/:repository/:oid", document.GetDocument)
	router.PUT("/data/:repository", document.UploadDocument)
	router.DELETE("/data/:repository/:oid", document.DeleteDocument)

	return router
}
