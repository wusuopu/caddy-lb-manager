package interfaces

type ICaddyfileService interface {
	GenCaddyfile() (string, error)
	ReloadFile(string) (bool, error)
	Reload(string) (bool, error)
	Validate(string) (bool, error)
	HashPassword (password string) (string, error)
}