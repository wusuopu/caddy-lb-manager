package models

import (
	"gorm.io/gorm"
)


type UpStream struct {
	gorm.Model
	Name				string						`gorm:"type:varchar(100);"`
	Scheme			string						`gorm:"type:varchar(20);"`		// eg: https://
	Address			string						`gorm:"type:varchar(300);"`		// host:port
}