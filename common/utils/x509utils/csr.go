package x509utils

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

func ReadCsrFile(csrFile string) (string, error) {
	data, err := os.ReadFile(csrFile)
	if err != nil {
		return "", err
	}

	b, _ := pem.Decode(data)
	if b == nil {
		_, err = x509.ParseCertificateRequest(data)
	} else {
		_, err = x509.ParseCertificateRequest(b.Bytes)
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}
