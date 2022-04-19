package models

import (
    "gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func NewDB(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    db.AutoMigrate(Plants{})

    return db, nil
}
