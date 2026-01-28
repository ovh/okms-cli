package restflags

import (
	"errors"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
)

type ProtectionLevel types.ProtectionLevelEnum

const (
	SOFTWARE   ProtectionLevel = ProtectionLevel(types.SOFTWARE)
	HSM        ProtectionLevel = ProtectionLevel(types.HSM)
	MANAGEDHSM ProtectionLevel = ProtectionLevel(types.MANAGEDHSM)
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
	case "managedHsm":
		*e = MANAGEDHSM
	default:
		return errors.New(`must be one of "soft", "hsm", "managedHsm"`)
	}
	return nil
}

func (e *ProtectionLevel) Type() string {
	return "soft|hsm|managedHsm"
}

func (e ProtectionLevel) RestModel() types.ProtectionLevelEnum {
	return types.ProtectionLevelEnum(e)
}
