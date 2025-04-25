package initialize

import (
	"app/di"
	"app/interfaces"
	"app/services"
)

func InitServices() {
	di.Service.CaddyfileService = interfaces.ICaddyfileService(new(services.CaddyfileService))
}