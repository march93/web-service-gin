package controllers

import (
	"errors"
	"net/http"
	"web-service-gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentController struct {
	DB *gorm.DB
}

// Get document by oid
func (d DocumentController) GetDocument(c *gin.Context) {
	name := c.Param("repository")
	oid := c.Param("oid")

	repository, err := models.GetSpecificDocument(d.DB, name, oid)
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

	c.JSON(http.StatusOK, documents[0])
}

// Upload a document
func (d DocumentController) UploadDocument(c *gin.Context) {
	name := c.Param("repository")
	var repository models.Repository
	var document models.Document

	// Bad input types provided
	err := c.BindJSON(&document)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	tx := d.DB.Begin()
	repository, err = models.GetSpecificDocumentByContent(tx, name, &document)
	if err != nil {
		// Repository does not exist, so we will create one
		// and attach our document to it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// First set the repository name to be used as primary key
			repository.Name = name
			repository, err = models.CreateRepository(tx, &repository)
			if err != nil {
				// Generic 500 server error
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
				return
			}

			// Create our document and set the repository name
			document.RepositoryName = name
			document, err = models.CreateDocument(tx, &document)
			if err != nil {
				// Generic 500 server error
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
				return
			}

			// Return oid and size object as part of 200 response
			tx.Commit()
			c.JSON(http.StatusCreated, gin.H{"oid": document.Oid, "size": len(document.Content)})
			return
		}

		// Generic 500 server error
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	// Repository already exists, check for a matching document
	if len(repository.Documents) == 0 {
		// Document with the same content does not exist, so we can create the document
		// Create our document and set the repository name
		document.RepositoryName = name
		document, err = models.CreateDocument(d.DB, &document)
		if err != nil {
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

	repository, err := models.GetSpecificDocument(d.DB, name, oid)
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
