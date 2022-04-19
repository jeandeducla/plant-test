package models

import "github.com/jinzhu/gorm"

type Asset struct {
    gorm.Model
    Name     string
    MaxPower uint
    Type     string
    PlantID  uint
}
