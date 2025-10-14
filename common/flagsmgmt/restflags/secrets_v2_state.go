package restflags

import (
	"errors"

	"github.com/ovh/okms-sdk-go/types"
)

type SecretV2State string

// Defines values for SecretV2State.
const (
	SecretV2StateActive      SecretV2State = "active"
	SecretV2StateDeactivated SecretV2State = "deactivated"
	SecretV2StateDeleted     SecretV2State = "deleted"
)

func (e *SecretV2State) String() string {
	return string(*e)
}

func (e *SecretV2State) Set(v string) error {
	switch v {
	case "active", "deactivated", "deleted":
		*e = SecretV2State(v)
		return nil
	default:
		return errors.New(`must be one of "active", "deactivated", "deleted"`)
	}
}

func (e *SecretV2State) ToRestSecretV2States() types.SecretV2State {
	return types.SecretV2State(*e)
}

func (e *SecretV2State) Type() string {
	return "active|deactivated|deleted"
}
