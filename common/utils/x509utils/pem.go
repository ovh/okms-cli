package x509utils

import (
	"encoding/pem"
	"errors"
)

func PemDecode(data []byte) (p *pem.Block, err error) {
	p, _ = pem.Decode(data)
	if p == nil {
		err = errors.New("Failed to decode CSR PEM block")
	}
	return
}
