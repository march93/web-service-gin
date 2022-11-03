package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	Oid       string `gorm:"primary_key;"`
	Content   string `gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (document *Document) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("Oid", uuid)
	return nil
}

// Fetch a document
func GetDocument(db *gorm.DB, document *Document, oid string) (err error) {
	err = db.Where("oid = ?", oid).First(document).Error

	if err != nil {
		return err
	}

	return nil
}

// Fetch a document by content
func GetDocumentByContent(db *gorm.DB, document *Document, newDocument *Document) (err error) {
	err = db.Where("content = ?", newDocument.Content).First(document).Error

	if err != nil {
		return err
	}

	return nil
}

// Upload a document
func CreateDocument(db *gorm.DB, document *Document) (err error) {
	err = db.Create(document).Error

	if err != nil {
		return err
	}

	return nil
}

// Update a document
func UpdateDocument(db *gorm.DB, document *Document) (err error) {
	err = db.Save(document).Error

	if err != nil {
		return err
	}

	return nil
}

// Delete a document
func DeleteDocument(db *gorm.DB, oid string) (err error) {
	err = db.Where("oid = ?", oid).Delete(&Document{}).Error

	if err != nil {
		return err
	}

	return nil
}
