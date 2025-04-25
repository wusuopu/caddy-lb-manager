package models

import (
	"fmt"
	"strings"

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

func (s *Server) GetAddress() string {
	var address []string
	if s.EnableSSL {
		address = append(address, "https://")
	}
	if s.Host != "" {
		address = append(address, s.Host)
	}
	address = append(address, fmt.Sprintf(":%d", s.Port))
	return strings.Join(address, "")
}