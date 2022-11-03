package controllers

import (
	"errors"
	"net/http"
	"web-service-gin/database"
	"web-service-gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentController struct {
	DB *gorm.DB
}

func NewDocumentController() *DocumentController {
	db := database.Db
	db.AutoMigrate(&models.Document{})
	return &DocumentController{DB: db}
}

// Get document by oid
func (d DocumentController) GetDocument(c *gin.Context) {
	oid := c.Param("oid")
	var document models.Document

	err := models.GetDocument(d.DB, &document, oid)
	if err != nil {
		// Throw 404 if document not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Generic 500 server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	// Return document with status 200
	c.JSON(http.StatusOK, document)
}
