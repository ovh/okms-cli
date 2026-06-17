package kmip

import (
	"github.com/ovh/okms-cli/common/utils/exit"
	explorer "github.com/phsym/kmip-explorer"
	"github.com/spf13/cobra"
)

func explorerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "explorer",
		Aliases: []string{"explore", "browse"},
		Short:   "Browse and manage kmip objects in an interactive terminal UI",
		Long: "Start the kmip-explorer terminal UI connected to the " +
			"configured KMIP endpoint. The UI takes over the terminal " +
			"and blocks until you quit it by pressing 'q'.",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// New(client, version, latestVersion); empty versions disables version display.
			exit.OnErr(explorer.New(kmipClient, "", "").Run())
		},
	}
}
