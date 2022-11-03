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

// Upload a document
func (d DocumentController) UploadDocument(c *gin.Context) {
	var existingDocument models.Document
	var newDocument models.Document
	c.BindJSON(&newDocument)

	err := models.GetDocumentByContent(d.DB, &existingDocument, &newDocument)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Document to be added doesn't have content matched to an existing document
		// so we can upload our new document directly to the database
		err = models.CreateDocument(d.DB, &newDocument)
		if err != nil {
			if errors.Is(err, gorm.ErrInvalidData) {
				// Passed in invalid data - return 400 bad request
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			// Generic 500 server error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}

		// Return oid and size object as part of 200 response
		c.JSON(http.StatusCreated, gin.H{"oid": newDocument.Oid, "size": len(newDocument.Content)})
		return
	}

	// Document not updated, no data returned
	c.Status(http.StatusNoContent)
}
