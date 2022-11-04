package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	Oid            string `gorm:"primary_key;"`
	Content        string `binding:"required,min=1" gorm:"not null;"`
	RepositoryName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (document *Document) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("Oid", uuid)
	return nil
}

// Create a document
func CreateDocument(db *gorm.DB, document *Document) (Document, error) {
	err := db.Create(&document).Error

	if err != nil {
		return *document, err
	}

	return *document, nil
}

// Delete a document
func DeleteDocument(db *gorm.DB, oid string) error {
	err := db.Where("oid = ?", oid).Delete(&Document{}).Error

	if err != nil {
		return err
	}

	return nil
}
