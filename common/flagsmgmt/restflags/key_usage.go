package restflags

import (
	"errors"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
)

type KeyUsage string

const (
	SIGN       KeyUsage = "sign"
	VERIFY     KeyUsage = "verify"
	ENCRYPT    KeyUsage = "encrypt"
	DECRYPT    KeyUsage = "decrypt"
	WRAPKEY    KeyUsage = "wrapKey"
	UNWRAPKEY  KeyUsage = "unwrapKey"
	DERIVEKEY  KeyUsage = "deriveKey"
	DERIVEBITS KeyUsage = "deriveBits"
)

func (e *KeyUsage) String() string {
	return string(*e)
}

func (e *KeyUsage) Set(v string) error {
	switch v {
	case "sign", "verify", "encrypt", "decrypt", "wrapKey", "unwrapKey", "deriveKey", "deriveBits":
		*e = KeyUsage(v)
		return nil
	default:
		return errors.New(`must be one of "sign", "verify", "encrypt", "decrypt", "wrapKey", "unwrapKey", "deriveKey", "deriveBits"`)
	}
}

func (e *KeyUsage) Type() string {
	return "sign|verify|encrypt|decrypt|wrapKey|unwrapKey|deriveKey|deriveBits"
}

type KeyUsageList []KeyUsage

func (e *KeyUsageList) Set(value string) error {
	arr := strings.Split(value, ",")
	// in case the user set the usage in the following form: --usage encrypt,decrypt:
	// split the input value, and trim spaces
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}
	for _, v := range arr {
		var ku KeyUsage
		if err := ku.Set(v); err != nil {
			return err
		}
		*e = append(*e, ku)
	}
	return nil
}

func (e *KeyUsageList) String() string {
	l := make([]string, 0)
	for _, ku := range *e {
		l = append(l, string(ku))
	}

	return strings.Join(l, ",")
}

func (e *KeyUsageList) ToStringArray() []string {
	l := make([]string, 0, len(*e))
	for _, ku := range *e {
		l = append(l, string(ku))
	}

	return l
}

func (e *KeyUsageList) ToCryptographicUsage() []types.CryptographicUsages {
	l := make([]types.CryptographicUsages, 0, len(*e))
	for _, ku := range *e {
		l = append(l, types.CryptographicUsages(ku))
	}

	return l
}

func (e *KeyUsageList) Type() string {
	return "Combination of: sign|verify|encrypt|decrypt|wrapKey|unwrapKey|deriveKey|deriveBits"
}
