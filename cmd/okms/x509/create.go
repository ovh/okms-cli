package x509

import "github.com/spf13/cobra"

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "gen", "generate"},
		Short:   "Generate certificates and CSR signed with a KMS key",
	}

	cmd.AddCommand(
		createGenerateCaCommand(),
		createGenerateCertCommand(),
		createGenerateCsrCommand(),
	)
	return cmd
}
