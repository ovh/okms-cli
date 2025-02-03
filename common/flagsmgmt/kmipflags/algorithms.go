package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type SymmetricAlg kmip.CryptographicAlgorithm

const (
	AES      SymmetricAlg = SymmetricAlg(kmip.AES)
	TDES     SymmetricAlg = SymmetricAlg(kmip.TDES)
	SKIPJACK SymmetricAlg = SymmetricAlg(kmip.SKIPJACK)
)

func (e *SymmetricAlg) String() string {
	return ttlv.EnumStr(kmip.CryptographicAlgorithm(*e))
}

func (e *SymmetricAlg) Set(v string) error {
	switch strings.ToLower(v) {
	case "aes":
		*e = AES
	case "tdes":
		*e = TDES
	case "skipjack":
		*e = SKIPJACK
	default:
		return errors.New(`must be one of "AES", "TDES", "SKIPJACK"`)
	}
	return nil
}

func (e *SymmetricAlg) Type() string {
	return "AES|TDES|SKIPJACK"
}

type AsymmetricAlg kmip.CryptographicAlgorithm

const (
	RSA   AsymmetricAlg = AsymmetricAlg(kmip.RSA)
	ECDSA AsymmetricAlg = AsymmetricAlg(kmip.ECDSA)
)

func (e *AsymmetricAlg) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.CryptographicAlgorithm(*e))
}

func (e *AsymmetricAlg) Set(v string) error {
	switch strings.ToLower(v) {
	case "rsa":
		*e = RSA
	case "ecdsa":
		*e = ECDSA
	default:
		return errors.New(`must be one of "RSA", "ECDSA"`)
	}
	return nil
}

func (e *AsymmetricAlg) Type() string {
	return "RSA|ECDSA"
}
