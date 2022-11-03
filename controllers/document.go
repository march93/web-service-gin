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
	db.AutoMigrate(&models.Repository{})
	return &DocumentController{DB: db}
}

// Get document by oid
func (d DocumentController) GetDocument(c *gin.Context) {
	name := c.Param("repository")
	oid := c.Param("oid")
	var repository models.Repository

	err := models.GetSpecificDocument(d.DB, &repository, name, oid)
	if err != nil {
		// Throw 404 if repository not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Generic 500 server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	documents := repository.Documents
	if len(documents) == 0 {
		// No document found within this repo
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, documents[0])

}

// Upload a document
func (d DocumentController) UploadDocument(c *gin.Context) {
	name := c.Param("repository")
	var repository models.Repository
	var document models.Document
	c.BindJSON(&document)

	err := models.GetSpecificDocumentByContent(d.DB, &repository, name, &document)
	if err != nil {
		// Repository does not exist, so we will create one
		// and attach our document to it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// First set the repository name to be used as primary key
			repository.Name = name
			err = models.CreateRepository(d.DB, &repository)
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

			// Create our document and set the repository name
			document.RepositoryName = name
			err = models.CreateDocument(d.DB, &document)
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
			c.JSON(http.StatusCreated, gin.H{"oid": document.Oid, "size": len(document.Content)})
			return
		}

		// Generic 500 server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	// Repository already exists, check for a matching document
	if len(repository.Documents) == 0 {
		// Document with the same content does not exist, so we can create the document
		// Create our document and set the repository name
		document.RepositoryName = name
		err = models.CreateDocument(d.DB, &document)
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
		c.JSON(http.StatusCreated, gin.H{"oid": document.Oid, "size": len(document.Content)})
		return
	}

	// Document not updated, no data returned
	c.Status(http.StatusNoContent)
}

// Delete a document
func (d DocumentController) DeleteDocument(c *gin.Context) {
	name := c.Param("repository")
	oid := c.Param("oid")
	var repository models.Repository

	err := models.GetSpecificDocument(d.DB, &repository, name, oid)
	if err != nil {
		// Throw 404 if repository not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "Repository not found"})
			return
		}

		// Generic 500 server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	documents := repository.Documents
	if len(documents) == 0 {
		// No document found within this repo
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "Document not found"})
		return
	}

	// Delete the document
	err = models.DeleteDocument(d.DB, documents[0].Oid)
	if err != nil {
		// Generic 500 server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.Status(http.StatusOK)
}
