package models

import "github.com/jinzhu/gorm"

type EnergyManager struct {
    gorm.Model
    Name    string
    Surname string
    Plants  []Plant `gorm:"constraint:OnDelete:SET NULL;"`
}
