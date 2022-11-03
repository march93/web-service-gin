package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	OID     string `gorm:"primarykey"`
	Content string
	Size    int64
}

// Fetch a document
func GetDocument(db *gorm.DB, document *Document, oid string) (err error) {
	err = db.Where("oid = ?", oid).First(document).Error

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
func DeleteDocument(db *gorm.DB, document *Document, oid string) (err error) {
	err = db.Where("oid = ?", oid).Delete(document).Error

	if err != nil {
		return err
	}

	return nil
}
