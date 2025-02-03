package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type State kmip.State

const (
	StatePreActive            State = State(kmip.StatePreActive)
	StateActive               State = State(kmip.StateActive)
	StateDeactivated          State = State(kmip.StateDeactivated)
	StateCompromised          State = State(kmip.StateCompromised)
	StateDestroyed            State = State(kmip.StateDestroyed)
	StateDestroyedCompromised State = State(kmip.StateDestroyedCompromised)
)

func (e *State) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.State(*e))
}

func (e *State) Set(v string) error {
	switch strings.ToLower(v) {
	case "preactive":
		*e = StatePreActive
	case "active":
		*e = StateActive
	case "deactivated":
		*e = StateDeactivated
	case "compromised":
		*e = StateCompromised
	case "destroyed":
		*e = StateDestroyed
	case "destroyedcompromised":
		*e = StateDestroyedCompromised
	default:
		return errors.New(`must be one of "PreActive", "Active", "Deactivated", "Compromised", "Destroyed", "DestroyedCompromised"`)
	}
	return nil
}

func (e *State) Type() string {
	return "PreActive|Active|Deactivated|Compromised|Destroyed|DestroyedCompromised"
}
