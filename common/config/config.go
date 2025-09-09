package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var k = koanf.New(".")

type EndpointConfig struct {
	Endpoint string
	CaFile   string
	Auth     EndpointAuth
}

func (epCfg EndpointConfig) TlsConfig(serverName string) *tls.Config {
	pool := loadCertPool(epCfg.CaFile)
	return &tls.Config{
		RootCAs:      pool,
		Certificates: epCfg.Auth.TlsCertificates(),
		ServerName:   serverName,
		MinVersion:   tls.VersionTLS12,
	}
}

func ReadUserInput(prompt, key, profile string, validate ...Validator) {
	// Build the key according to the given profile
	profiledKey := fmt.Sprintf("profiles.%s.%s", profile, key)
	curVal := k.String(profiledKey)

	// Format the prompt: '<prompt> [<current value>]: '
	// ie. 'CA file [path/to/ca.crt]: '
	var userVal string

	for {
		userVal = exit.OnErr2(pterm.DefaultInteractiveTextInput.
			WithDefaultValue(curVal).
			WithOnInterruptFunc(func() {
				exit.OnErr(errors.New("Interrupted"))
			}).
			Show(prompt))
		if err := checkValidate(userVal, validate...); err != nil {
			fmt.Fprintln(os.Stderr, "Invalid value:", err.Error())
			continue
		}
		break
	}
	exit.OnErr(k.Set(profiledKey, strings.TrimSpace(userVal)))
}

func SetConfigKey(profile, key, value string) {
	// Build the key according to the given profile
	profiledKey := fmt.Sprintf("profiles.%s.%s", profile, key)
	exit.OnErr(k.Set(profiledKey, strings.TrimSpace(value)))
}

func LoadFromFile(defaultFile, customFile string) (string, error) {
	k = koanf.New(".")
	defer checkVersionAndMigrate()
	if customFile == "" {
		defaultFile += ".yaml"
		if err := k.Load(file.Provider(defaultFile), yaml.Parser()); err != nil {
			home, _ := os.UserHomeDir()
			fp := filepath.Join(home, ".ovh-kms", defaultFile)
			if err = k.Load(file.Provider(fp), yaml.Parser()); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Creating empty config file at %s\n", fp)
				exit.OnErr(k.Set("version", 1))
				_ = WriteToFile(fp)
			}
			return fp, nil
		}
		return defaultFile, nil
	}

	fpath, err := utils.ExpandTilde(customFile)
	if err != nil {
		return "", err
	}
	// While we ignore the errors while reading default (and implicit) config file,
	// errors on reading explicitly set config file are returned.
	return fpath, k.Load(file.Provider(fpath), yaml.Parser())
}

func LoadEndpointConfig(command *cobra.Command, service, configFile string) EndpointConfig {
	defaultFile := "okms"

	if _, err := LoadFromFile(defaultFile, configFile); err != nil {
		exit.OnErr(fmt.Errorf("Failed to load config file: %w", err))
	}
	return loadV1(command, service)
}

func ProfileList() []string {
	list := k.MapKeys("profiles")
	if len(list) == 0 {
		list = []string{"default"}
	}
	return list
}

func CurrentProfile() string {
	current := k.String("profile")
	if current == "" {
		current = ProfileList()[0]
	}
	return current
}

func SwitchProfile(profile string) error {
	if !k.Exists("profiles." + profile) {
		return fmt.Errorf("Profile %q is not defined", profile)
	}
	return k.Set("profile", profile)
}

func loadV1(command *cobra.Command, service string) EndpointConfig {
	if service == "" {
		panic("Service name cannot be empty")
	}
	if k.Int("version") != 1 {
		panic("Invalid config version (expect 1)")
	}

	profile := GetString(k, "profile", "KMS_PROFILE", command.Flags().Lookup("profile"))
	if profile == "" {
		profile = "default"
	}

	envPrefix := "KMS_" + strings.ToUpper(service)
	k := k.Cut(fmt.Sprintf("profiles.%s.%s", profile, service))

	ep := GetString(k, "endpoint", envPrefix+"_ENDPOINT", command.Flags().Lookup("endpoint"))
	if ep == "" {
		exit.OnErr(errors.New("Missing endpoint address parameter"))
	}
	caFile := GetString(k, "ca", envPrefix+"_CA", command.Flags().Lookup("ca"))

	authMethod := GetString(k, "auth.type", envPrefix+"_AUTH_METHOD", command.Flags().Lookup("auth-method"))

	return EndpointConfig{
		Endpoint: ep,
		CaFile:   exit.OnErr2(utils.ExpandTilde(caFile)),
		Auth:     buildAuthMethod(service, authMethod, command, k.Cut("auth"), envPrefix),
	}
}

func GetString(k *koanf.Koanf, key, env string, flag *pflag.Flag) string {
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}
	if env != "" {
		if envVal, ok := os.LookupEnv(env); ok {
			return envVal
		}
	}
	if key != "" && k.Exists(key) {
		return k.String(key)
	}
	if flag != nil && !flag.Changed {
		return flag.DefValue
	}
	return ""
}

func WriteToFile(customFile string) error {
	var err error
	bytes, err := k.Marshal(yaml.Parser())
	if err != nil {
		return err
	}
	if customFile == "" {
		return errors.New("No config file path given")
	}
	if customFile, err = utils.ExpandTilde(customFile); err != nil {
		return err
	}
	cfgPath := filepath.Dir(customFile)

	if err = os.MkdirAll(cfgPath, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create folder:", cfgPath)
		return err
	}

	if err = os.WriteFile(customFile, bytes, 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to save configuration file:", err.Error())
		return err
	}

	return nil
}

func checkVersionAndMigrate() {
	switch k.Int("version") {
	case 0:
		if len(k.MapKeys("")) > 0 {
			fmt.Fprintf(os.Stderr, "[WARN] Using an old configuration format (v0) Migrate by running %s configure --migrate-only\n", os.Args[0])
		}
		k = migrateV0toV1(k)
		return
	case 1:
		return
	default:
		exit.OnErr(errors.New("Unsupported config version"))
	}
}

func migrateV0toV1(k *koanf.Koanf) *koanf.Koanf {
	if version := k.Int("version"); version == 1 {
		return k
	} else if version != 0 {
		panic("Cannot migrate to config v1: Invalid config version (must be 0)")
	}

	nk := koanf.New(".")
	exit.OnErr(nk.Set("version", 1))
	exit.OnErr(nk.Set("profile", "default"))

	for _, key := range k.MapKeys("") {
		prefix := "profiles." + key + "."
		pk := k.Cut(key)

		migrateServiceConfigV0toV1(pk, nk, prefix+"http.")

		for _, svc := range pk.MapKeys("") {
			pk := pk.Cut(svc)
			targetPrefix := prefix + svc + "."
			migrateServiceConfigV0toV1(pk, nk, targetPrefix)
		}
	}
	return nk
}

func migrateServiceConfigV0toV1(src, dst *koanf.Koanf, targetPrefix string) {
	if ep := src.String("endpoint"); ep != "" {
		exit.OnErr(dst.Set(targetPrefix+"endpoint", ep))
	}
	if ca := src.String("ca"); ca != "" {
		exit.OnErr(dst.Set(targetPrefix+"ca", ca))
	}
	if cert := src.String("cert"); cert != "" {
		exit.OnErr(dst.Set(targetPrefix+"auth.type", "mtls"))
		exit.OnErr(dst.Set(targetPrefix+"auth.cert", cert))
	}
	if key := src.String("key"); key != "" {
		exit.OnErr(dst.Set(targetPrefix+"auth.type", "mtls"))
		exit.OnErr(dst.Set(targetPrefix+"auth.key", key))
	}
}
