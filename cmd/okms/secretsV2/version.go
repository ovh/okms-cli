package secretsv2

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/restflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
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
	)
	return cmd
}

func secretVersionGetCmd() *cobra.Command {
	var (
		version     uint32
		includeData bool
	)
	cmd := &cobra.Command{
		Use:   "get PATH --version VERSION ",
		Short: "Retrieve a secret version",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			var v uint32
			if cmd.Flag("version").Changed {
				v = version
			}

			resp := exit.OnErr2(common.Client().GetSecretVersionV2(cmd.Context(), args[0], v, &includeData))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				createdAt := resp.CreatedAt
				deactivatedAt := utils.DerefOrDefault(resp.DeactivatedAt)
				id := resp.Id
				state := resp.State

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Created at", createdAt},
					{"Deactivated at", deactivatedAt},
					{"Id ", fmt.Sprintf("%b", id)},
					{"State", string(state)},
				})
				table.Render()

				if includeData && resp.Data != nil {
					fmt.Println("Data")
					tableData := tablewriter.NewWriter(os.Stdout)
					tableData.SetHeader([]string{"Key", "Value"})
					for k, v := range *resp.Data {
						tableData.Append([]string{k, fmt.Sprintf("%v", v)})
					}
					tableData.Render()
				}
			}
		},
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version. If not set, the latest version will be returned.")
	cmd.Flags().BoolVar(&includeData, "include-data", true, "Include the secret data. If not set they will be returned.")
	return cmd
}

// TODO : think of a better way to display the list of versions and metadata.
func secretVersionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get PATH --version VERSION ",
		Short: "Retrieve a secret version",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			resp := exit.OnErr2(common.Client().ListSecretVersionV2(cmd.Context(), args[0]))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				for _, version := range *resp {

					createdAt := version.CreatedAt
					deactivatedAt := utils.DerefOrDefault(version.DeactivatedAt)
					id := version.Id
					state := version.State

					fmt.Println("Metadata")
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Key", "Value"})
					table.AppendBulk([][]string{
						{"Created at", createdAt},
						{"Deactivated at", deactivatedAt},
						{"Id ", fmt.Sprintf("%b", id)},
						{"State", string(state)},
					})
					table.Render()

					if version.Data != nil {
						fmt.Println("Data")
						tableData := tablewriter.NewWriter(os.Stdout)
						tableData.SetHeader([]string{"Key", "Value"})
						for k, v := range *version.Data {
							tableData.Append([]string{k, fmt.Sprintf("%v", v)})
						}
						tableData.Render()
					}
				}
			}
		},
	}
	return cmd
}

// TODO : Should we provide an easier command like
// activate, deactivate or delete directly with the version and the path ?
func secretVersionPutCmd() *cobra.Command {
	var (
		version uint32
		state   restflags.SecretV2State
	)
	cmd := &cobra.Command{
		Use:   "update  PATH --version VERSION --state STATE",
		Short: "Update a secret version",
		Args:  cobra.MinimumNArgs(4),
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

			resp := exit.OnErr2(common.Client().PutSecretVersionV2(cmd.Context(), args[0], version, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				createdAt := resp.CreatedAt
				deactivatedAt := utils.DerefOrDefault(resp.DeactivatedAt)
				id := resp.Id
				state := resp.State

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Created at", createdAt},
					{"Deactivated at", deactivatedAt},
					{"Id ", fmt.Sprintf("%b", id)},
					{"State", string(state)},
				})
				table.Render()
			}
		},
	}
	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version. If not set, the latest version will be returned.")
	cmd.Flags().Var(&state, "state", "State of the secret version. Value must be in active|deactivated|deleted")
	return cmd
}

func secretVersionPostCmd() *cobra.Command {
	var (
		cas uint32
		// TODO Add customMetadata ? How ?
	)
	cmd := &cobra.Command{
		Use:   "create [FLAGS] PATH [DATA]",
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

			resp := exit.OnErr2(common.Client().PostSecretVersionV2(cmd.Context(), args[0], c, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				createdAt := resp.CreatedAt
				deactivatedAt := utils.DerefOrDefault(resp.DeactivatedAt)
				id := resp.Id
				state := resp.State

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Created at", createdAt},
					{"Deactivated at", deactivatedAt},
					{"Id ", fmt.Sprintf("%b", id)},
					{"State", string(state)},
				})
				table.Render()
			}
		},
	}
	cmd.Flags().Uint32Var(&cas, "cas", 0, "Secret version number. Required if cas-required is set to true.")

	return cmd
}
