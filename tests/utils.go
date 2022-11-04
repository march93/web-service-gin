package tests

import "gorm.io/gorm"

func closeDB(db *gorm.DB) {
	dbInstance, _ := db.DB()
	_ = dbInstance.Close()
}
