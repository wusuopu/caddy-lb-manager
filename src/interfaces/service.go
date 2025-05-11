package interfaces

import "app/types"

type ICaddyfileService interface {
	GenCaddyfile() (string, error)
	ReloadFile(string) (bool, error)
	Reload(string) (bool, error)
	Validate(string) (bool, error)
	TouchReloadTime() (bool, error)
	PullConfigAndReload() (bool, error)
	HashPassword (password string) (string, error)
	ListCertificate() []types.CertificateInfo
}