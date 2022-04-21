package models

import "github.com/jinzhu/gorm"

type Plant struct {
    gorm.Model
    Name            string
    Address         string
    MaxPower        uint
    EnergyManagerID uint
    Assets          []Asset `gorm:"constraint:OnDelete:CASCADE;"`
}

