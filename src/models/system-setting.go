package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type SystemSetting struct {
	gorm.Model
	Name				string						`gorm:"type:varchar(100);"`
	Value				datatypes.JSON		`gorm:"type:json;"`
}