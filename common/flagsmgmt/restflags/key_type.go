package restflags

import (
	"errors"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
)

type KeyType types.KeyTypes

const (
	OCTETSTREAM    KeyType = KeyType(types.Oct)
	RSA            KeyType = KeyType(types.RSA)
	ELLIPTIC_CURVE KeyType = KeyType(types.EC)
)

func (e *KeyType) String() string {
	return string(*e)
}

func (e *KeyType) Set(v string) error {
	switch strings.ToLower(v) {
	case "oct":
		*e = OCTETSTREAM
	case "rsa":
		*e = RSA
	case "ec":
		*e = ELLIPTIC_CURVE
	default:
		return errors.New(`must be one of "oct", "rsa", "ec"`)
	}
	return nil
}

func (e *KeyType) Type() string {
	return "oct|rsa|ec"
}

func (e KeyType) RestModel() types.KeyTypes {
	return types.KeyTypes(e)
}
