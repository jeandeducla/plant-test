package models

import "github.com/jinzhu/gorm"

type Plants struct {
    gorm.Model
    Name string
    Address string
    MaxPower uint
}

