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
func CreateRepository(db *gorm.DB, repository *Repository) (Repository, error) {
	err := db.Create(&repository).Error

	if err != nil {
		return *repository, err
	}

	return *repository, nil
}

// Fetch document from a repository by oid
func GetSpecificDocument(db *gorm.DB, name string, oid string) (Repository, error) {
	var repository Repository
	err := db.Where("name = ?", name).Preload("Documents", "oid = ?", oid).First(&repository).Error

	if err != nil {
		return repository, err
	}

	return repository, nil
}

// Fetch document from a repository by content
func GetSpecificDocumentByContent(db *gorm.DB, name string, document *Document) (Repository, error) {
	var repository Repository
	err := db.Where("name = ?", name).Preload("Documents", "content = ?", document.Content).First(&repository).Error

	if err != nil {
		return repository, err
	}

	return repository, nil
}
