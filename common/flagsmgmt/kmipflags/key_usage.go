package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type KeyUsage kmip.CryptographicUsageMask

const (
	SIGN      KeyUsage = KeyUsage(kmip.Sign)
	VERIFY    KeyUsage = KeyUsage(kmip.Verify)
	ENCRYPT   KeyUsage = KeyUsage(kmip.Encrypt)
	DECRYPT   KeyUsage = KeyUsage(kmip.Decrypt)
	WRAPKEY   KeyUsage = KeyUsage(kmip.WrapKey)
	UNWRAPKEY KeyUsage = KeyUsage(kmip.UnwrapKey)
	DERIVEKEY KeyUsage = KeyUsage(kmip.DeriveKey)
)

func (e *KeyUsage) String() string {
	return ttlv.BitmaskStr(kmip.CryptographicUsageMask(*e), ",")
	// return "TODO"
}

func (e *KeyUsage) Set(v string) error {
	switch strings.ToLower(v) {
	case "sign":
		*e = SIGN
	case "verify":
		*e = VERIFY
	case "encrypt":
		*e = ENCRYPT
	case "decrypt":
		*e = DECRYPT
	case "wrapkey":
		*e = WRAPKEY
	case "unwrapkey":
		*e = UNWRAPKEY
	case "derivekey":
		*e = DERIVEKEY
	default:
		return errors.New(`must be one of "Sign", "Verify", "Encrypt", "Decrypt", "WrapKey", "UnwrapKey", "DeriveKey"`)
	}
	return nil
}

func (e *KeyUsage) Type() string {
	return "Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey"
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
		l = append(l, ku.String())
	}

	return strings.Join(l, ",")
}

// func (e *KeyUsageList) ToStringArray() []string {
// 	l := make([]string, 0, len(*e))
// 	for _, ku := range *e {
// 		l = append(l, ku.String())
// 	}

// 	return l
// }

func (e *KeyUsageList) ToCryptographicUsageMask() kmip.CryptographicUsageMask {
	var mask kmip.CryptographicUsageMask
	for _, ku := range *e {
		mask |= kmip.CryptographicUsageMask(ku)
	}

	return mask
}

func (e *KeyUsageList) Type() string {
	return "Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey"
}
