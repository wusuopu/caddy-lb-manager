package di

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	DB *gorm.DB
	Logger *zap.Logger
}

var Container = new(container)
