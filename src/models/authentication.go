package models

import (
	"gorm.io/gorm"
)


type Authentication struct {
	gorm.Model
	Name				string						`gorm:"type:varchar(100);"`
	Username		string						`gorm:"type:varchar(100);"`
	Password		string						`gorm:"type:varchar(100);"`
	HashedPw		string						`gorm:"type:varchar(100);"`
}