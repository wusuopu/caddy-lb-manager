package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type Route struct {
	gorm.Model
	Name				string						`gorm:"type:varchar(100);"`
	Methods			string						`gorm:"type:varchar(100);"`		// GET,POST,....
	Path				string						`gorm:"type:varchar(300);"`
	Header			datatypes.JSON			// {[field]: {value: string, isReg: bool}}
	StripPath		bool
	UpStreamId	uint
	UpStream		UpStream
	Enable			bool
	ServerId		uint
	AuthenticationId	uint
	Authentication Authentication
}