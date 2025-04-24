package models

import (
	"gorm.io/gorm"
)


type Server struct {
	gorm.Model
	Name				string						`gorm:"type:varchar(100);"`
	Host				string						`gorm:"type:varchar(100);"`
	Port				uint
	EnableSSL		bool
	Enable			bool
}