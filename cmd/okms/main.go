package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ovh/okms-cli/cmd/okms/configure"
	"github.com/ovh/okms-cli/cmd/okms/keys"
	"github.com/ovh/okms-cli/cmd/okms/kmip"

	"github.com/ovh/okms-cli/cmd/okms/x509"
	"github.com/ovh/okms-cli/common/commands"
	"github.com/ovh/okms-cli/common/config"

	"github.com/spf13/cobra"
)

var (
	// The following values will be set automatically by goreleaser during the CI/CD pipeline execution
	// see: https://goreleaser.com/cookbooks/using-main.version/ and https://goreleaser.com/customization/builds/
	// The default ldflags are '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}'.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func createRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               filepath.Base(os.Args[0]),
		DisableAutoGenTag: true, // Do not add timestamp in generated markdown to avoid useles diffs
	}

	config.SetupConfigFlags(command)

	command.AddCommand(
		// rnd.CreateCommand(nil),
		keys.CreateCommand(nil),
		// secrets.CreateCommand(nil),
		x509.CreateX509Command(nil),
		kmip.NewCommand(nil),
		configure.CreateCommand(),
		commands.NewMarkdownCmd(command),
		commands.NewVersionCmd(&version, &commit, &date),
	)

	return command
}

func main() {
	if err := createRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
