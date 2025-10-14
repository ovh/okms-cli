package secretsv2

import (
	"fmt"
	"io"
	"os"

	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/restflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

func secretVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "This command has subcommands for interacting with your secret's versions.",
	}

	cmd.AddCommand(
		secretVersionGetCmd(),
		secretVersionPutCmd(),
		secretVersionPostCmd(),
		secretVersionListCmd(),
		secretVersionActiveCmd(),
		secretVersionDeactivateCmd(),
		secretVersionDeleteCmd(),
	)
	return cmd
}

func secretVersionGetCmd() *cobra.Command {
	var (
		version     uint32
		includeData bool
	)
	cmd := &cobra.Command{
		Use:   "get PATH --version=VERSION ",
		Short: "Retrieve a secret version",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var v uint32
			if cmd.Flag("version").Changed {
				v = version
			}

			resp := exit.OnErr2(common.Client().GetSecretVersionV2(cmd.Context(), common.GetOkmsId(), args[0], v, &includeData))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadataVersion(*resp)
				if includeData && resp.Data != nil {
					renderDataVersion(*resp.Data)
				}
			}
		},
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version.")
	cmd.Flags().BoolVar(&includeData, "include-data", false, "Include the secret data. If not set they will not be returned.")
	return cmd
}

func secretVersionListCmd() *cobra.Command {
	var (
		pageSize uint32
	)
	cmd := &cobra.Command{
		Use:   "list PATH",
		Short: "Retrieve all secret versions",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			versions := types.ListSecretVersionV2Response{}

			for sec, err := range common.Client().ListAllSecretVersions(common.GetOkmsId(), args[0], &pageSize).Iter(cmd.Context()) {
				exit.OnErr(err)
				versions = append(versions, *sec)
			}

			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(versions)
			} else {
				renderListMetadataVersion(versions)
			}
		},
	}

	cmd.Flags().Uint32Var(&pageSize, "page-size", 100, "Number of secret versions to fetch per page (between 10 and 500)")
	return cmd
}

func secretVersionPutCmd() *cobra.Command {
	var (
		version uint32
		state   restflags.SecretV2State
	)
	cmd := &cobra.Command{
		Use:   "update  PATH --version VERSION --state STATE",
		Short: "Update a secret version",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flag("version").Changed {
				fmt.Fprintln(os.Stderr, "Missing flag version")
				os.Exit(1)
			}
			if !cmd.Flag("state").Changed {
				fmt.Fprintln(os.Stderr, "Missing flag state")
				os.Exit(1)
			}

			body := types.PutSecretVersionV2Request{
				State: state.ToRestSecretV2States(),
			}

			resp := exit.OnErr2(common.Client().PutSecretVersionV2(cmd.Context(), common.GetOkmsId(), args[0], version, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadataVersion(*resp)
			}
		},
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version.")
	cmd.Flags().Var(&state, "state", "State of the secret version. Value must be in active|deactivated|deleted")
	return cmd
}

func stateFunction(version *uint32, state types.SecretV2State) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("version").Changed {
			fmt.Fprintln(os.Stderr, "Missing flag version")
			os.Exit(1)
		}
		body := types.PutSecretVersionV2Request{
			State: state,
		}

		resp := exit.OnErr2(common.Client().PutSecretVersionV2(cmd.Context(), common.GetOkmsId(), args[0], *version, body))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			renderMetadataVersion(*resp)
		}
	}
}

func secretVersionActiveCmd() *cobra.Command {
	var (
		version uint32
	)
	cmd := &cobra.Command{
		Use:   "activate  PATH --version VERSION ",
		Short: "Activate a secret version",
		Args:  cobra.MinimumNArgs(1),
		Run:   stateFunction(&version, types.SecretV2StateActive),
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version.")
	return cmd
}

func secretVersionDeactivateCmd() *cobra.Command {
	var (
		version uint32
	)
	cmd := &cobra.Command{
		Use:   "deactivate  PATH --version VERSION ",
		Short: "Deactivate a secret version",
		Args:  cobra.MinimumNArgs(1),
		Run:   stateFunction(&version, types.SecretV2StateDeactivated),
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version. If not set, the latest version will be returned.")
	return cmd
}

func secretVersionDeleteCmd() *cobra.Command {
	var (
		version uint32
	)
	cmd := &cobra.Command{
		Use:   "delete  PATH --version VERSION ",
		Short: "Delete a secret version",
		Args:  cobra.MinimumNArgs(1),
		Run:   stateFunction(&version, types.SecretV2StateDeleted),
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version.")
	return cmd
}

func secretVersionPostCmd() *cobra.Command {
	var (
		cas uint32
	)
	cmd := &cobra.Command{
		Use:   "create PATH [DATA]",
		Short: "Create a secret version",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			body := types.PostSecretVersionV2Request{}
			var c *uint32
			if cmd.Flag("cas").Changed {
				c = &cas
			}

			in := io.Reader(os.Stdin)
			data, err := restflags.ParseArgsData(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}
			body.Data = &data

			resp := exit.OnErr2(common.Client().PostSecretVersionV2(cmd.Context(), common.GetOkmsId(), args[0], c, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadataVersion(*resp)
			}
		},
	}
	cmd.Flags().Uint32Var(&cas, "cas", 0, "Secret version number. Required if cas-required is set to true.")

	return cmd
}
