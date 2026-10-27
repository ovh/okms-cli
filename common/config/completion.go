package config

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/internal/utils"
	"github.com/spf13/cobra"
)

// CompleteProfiles is a completion function suggesting the profiles defined in
// the configuration file.
func CompleteProfiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	configFile, _ := cmd.Flags().GetString("config")
	profiles := profileListFromFile(configFile)
	if len(profiles) == 0 {
		profiles = []string{"default"}
	}
	return profiles, cobra.ShellCompDirectiveNoFileComp
}

// HasEndpointConfig reports whether an endpoint configuration may be available
// for the given service, either from the command line, the environment, or a
// configuration file. Shell completion functions use it to avoid loading (and
// possibly creating) an empty configuration file when the CLI has not been
// configured yet.
func HasEndpointConfig(cmd *cobra.Command, service string) bool {
	if flag := cmd.Flags().Lookup("endpoint"); flag != nil && flag.Changed {
		return true
	}

	envPrefixes := []string{"KMS_" + strings.ToUpper(service)}
	if strings.EqualFold(service, "restapi") {
		// KMS_HTTP_* is a transparent alias for KMS_RESTAPI_*.
		envPrefixes = append(envPrefixes, "KMS_HTTP")
	}
	for _, prefix := range envPrefixes {
		if os.Getenv(prefix+"_ENDPOINT") != "" {
			return true
		}
	}

	configFile, _ := cmd.Flags().GetString("config")
	for _, path := range configFilePaths(configFile) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

// profileListFromFile reads the profile names from the configuration file.
// Unlike LoadFromFile, it neither creates nor migrates the configuration file,
// as completion must not have side effects.
func profileListFromFile(configFile string) []string {
	for _, path := range configFilePaths(configFile) {
		cfg := koanf.New(".")
		if err := cfg.Load(file.Provider(path), yaml.Parser()); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil
		}

		if cfg.Int("version") == 0 {
			cfg = migrateV0toV1(cfg)
		}

		profiles := cfg.MapKeys("profiles")
		slices.Sort(profiles)
		return profiles
	}
	return nil
}

// configFilePaths returns the candidate paths of the configuration file, in
// lookup order, mirroring the ones used by LoadFromFile.
func configFilePaths(configFile string) []string {
	if configFile != "" {
		fpath, err := utils.ExpandTilde(configFile)
		if err != nil {
			return nil
		}
		return []string{fpath}
	}
	home, _ := os.UserHomeDir()
	return []string{
		"okms.yaml",
		filepath.Join(home, ".ovh-kms", "okms.yaml"),
	}
}
