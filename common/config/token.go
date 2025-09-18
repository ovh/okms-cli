package config

import (
	"crypto/tls"
	"errors"

	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func init() {
	RegisterAuthMethod("token", newTokenAuth)
}

func newTokenAuth(cmd *cobra.Command, k *koanf.Koanf, envPrefix string) EndpointAuth {
	token := GetString(k, "token", envPrefix+"_TOKEN", cmd.Flags().Lookup("token"))
	if token == "" {
		exit.OnErr(errors.New("Missing token parameter"))
	}

	okmsIdStr := GetString(k, "okmsId", envPrefix+"_OKMSID", cmd.Flags().Lookup("okmsId"))
	if okmsIdStr == "" {
		exit.OnErr(errors.New("Missing okmsId parameter"))
	}

	okmsId, err := uuid.Parse(okmsIdStr)
	if err != nil {
		exit.OnErr(errors.New("Invalid okmsId"))
	}

	return &tokenAuth{
		token:  token,
		okmsId: okmsId,
	}
}

type tokenAuth struct {
	token  string
	okmsId uuid.UUID
}

func (epCfg *tokenAuth) GetToken() *string {
	return &epCfg.token
}

func (epCfg *tokenAuth) GetOkmsId() uuid.UUID {
	return epCfg.okmsId
}

func (epCfg *tokenAuth) TlsCertificates() []tls.Certificate {
	return nil
}
