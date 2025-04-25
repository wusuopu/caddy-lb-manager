package types


type CertificateInfo struct {
	Name string
	Sans []string
	ValidityStart string
	ValidityEnd string
}