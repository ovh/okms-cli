package restflags

import (
	"errors"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
)

type ProtectionLevel types.ProtectionLevelEnum

const (
	SOFTWARE ProtectionLevel = ProtectionLevel(types.SOFTWARE)
	HSM      ProtectionLevel = ProtectionLevel(types.HSM)
)

func (e *ProtectionLevel) String() string {
	return string(*e)
}

func (e *ProtectionLevel) Set(v string) error {
	switch strings.ToLower(v) {
	case "soft":
		*e = SOFTWARE
	case "hsm":
		*e = HSM
	default:
		return errors.New(`must be one of "soft", "hsm"`)
	}
	return nil
}

func (e *ProtectionLevel) Type() string {
	return "soft|hsm"
}

func (e ProtectionLevel) RestModel() types.ProtectionLevelEnum {
	return types.ProtectionLevelEnum(e)
}
