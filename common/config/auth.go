package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

type EndpointAuth interface {
	TlsCertificates() []tls.Certificate
}

type AuthMethod func(*cobra.Command, *koanf.Koanf, string) EndpointAuth

var authMethods = map[string]AuthMethod{}

func SupportedAuthMethods() []string {
	names := make([]string, 0, len(authMethods))
	for k := range authMethods {
		if k == "" {
			continue
		}
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}

func RegisterAuthMethod(name string, method AuthMethod) {
	name = strings.ToLower(name)
	if _, ok := authMethods[name]; ok {
		panic(fmt.Sprintf("Auth method %s already registered", name))
	}
	if method == nil {
		panic("provided auth method cannot be null")
	}
	authMethods[name] = method
}

func getAuthMethod(name string) AuthMethod {
	name = strings.ToLower(name)
	authMethod, ok := authMethods[name]
	if !ok {
		exit.OnErr(errors.New("Unsupported auth type"))
	}
	return authMethod
}

func buildAuthMethod(name string, cmd *cobra.Command, k *koanf.Koanf, envPrefix string) EndpointAuth {
	return getAuthMethod(name)(cmd, k, envPrefix)
}

type AuthMethodFlag string

func (e *AuthMethodFlag) String() string {
	return string(*e)
}

func (e *AuthMethodFlag) Set(v string) error {
	for k := range authMethods {
		if k == v {
			*e = AuthMethodFlag(v)
			return nil
		}
	}
	return fmt.Errorf(`must be one of %s`, strings.Join(SupportedAuthMethods(), ", "))
}

func (e *AuthMethodFlag) Type() string {
	return strings.Join(SupportedAuthMethods(), "|")
}
