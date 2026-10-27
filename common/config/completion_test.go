package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteProfiles(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "okms.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`version: 1
profiles:
  staging: {}
  default: {}
`), 0o600))

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("config", "", "")
	require.NoError(t, cmd.Flags().Set("config", configFile))

	completions, directive := CompleteProfiles(cmd, nil, "")

	assert.Equal(t, []string{"default", "staging"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestCompleteProfilesWithoutConfigFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cmd := &cobra.Command{Use: "test"}

	completions, directive := CompleteProfiles(cmd, nil, "")

	assert.Equal(t, []string{"default"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestHasEndpointConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("KMS_RESTAPI_ENDPOINT", "")
	t.Setenv("KMS_HTTP_ENDPOINT", "")
	t.Setenv("KMS_KMIP_ENDPOINT", "")

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("endpoint", "", "")
	cmd.Flags().String("config", "", "")

	assert.False(t, HasEndpointConfig(cmd, "restapi"))

	t.Setenv("KMS_RESTAPI_ENDPOINT", "https://kms.example.com")
	assert.True(t, HasEndpointConfig(cmd, "restapi"))
	assert.False(t, HasEndpointConfig(cmd, "kmip"))
}
