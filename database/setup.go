package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

	"github.com/niranjandahal/hostedgo/models"

)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    db.AutoMigrate(&models.Item{})
    return db, nil
}
