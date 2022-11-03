package models

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	Name      string     `gorm:"primary_key;unique"`
	Documents []Document `gorm:"foreignKey:RepositoryName;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Create a repository
func CreateRepository(db *gorm.DB, repository *Repository) (err error) {
	err = db.Create(repository).Error

	if err != nil {
		return err
	}

	return nil
}

// Fetch document from a repository by oid
func GetSpecificDocument(db *gorm.DB, repository *Repository, name string, oid string) (err error) {
	err = db.Where("name = ?", name).Preload("Documents", "oid = ?", oid).First(repository).Error

	if err != nil {
		return err
	}

	return nil
}

// Fetch document from a repository by content
func GetSpecificDocumentByContent(db *gorm.DB, repository *Repository, name string, document *Document) (err error) {
	err = db.Where("name = ?", name).Preload("Documents", "content = ?", document.Content).First(repository).Error

	if err != nil {
		return err
	}

	return nil
}
