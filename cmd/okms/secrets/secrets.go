package secrets

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

func kvGetCmd() *cobra.Command {
	var (
		version uint32
	)

	cmd := &cobra.Command{
		Use:   "get PATH",
		Short: "Retrieves the value from KMS's key-value store at the given key name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var v *uint32
			if version != 0 {
				v = &version
			}

			resp := exit.OnErr2(common.Client().GetSecretRequest(cmd.Context(), args[0], v))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else if resp.Data != nil {
				renderSecretMetadataTable(resp.Data.Metadata)

				if resp.Data.Data != nil {
					fmt.Println("Data")
					table := tablewriter.NewWriter(os.Stdout)
					table.Header([]string{"Key", "Value"})
					kvs, ok := (*resp.Data.Data).(map[string]any)
					if ok {
						for k, v := range kvs {
							exit.OnErr(table.Append([]string{k, fmt.Sprintf("%v", v)}))
						}
					}
					exit.OnErr(table.Render())
				}
			}
		},
	}

	cmd.Flags().Uint32Var(&version, "version", 0, "If passed, the value at the version number will be returned")
	return cmd
}

func kvPutCmd() *cobra.Command {
	var (
		cas int32
	)

	cmd := &cobra.Command{
		Use:   "put PATH [DATA]",
		Short: "Writes the data to the given path in the key-value store. (DATA format: bar=baz)",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)

			data, err := restflags.ParseArgsData(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}

			var c uint32
			if cas != -1 {
				c = utils.ToUint32(c)
			}
			body := types.PostSecretRequest{
				Data: new(any),
				Options: &types.PostSecretOptions{
					Cas: &c,
				},
			}

			*(body.Data) = data

			resp := exit.OnErr2(common.Client().PostSecretRequest(cmd.Context(), args[0], body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderSecretMetadataTable(resp.Data)
			}
		},
	}

	cmd.Flags().Int32Var(&cas, "cas", -1, "Specifies to use a Check-And-Set operation. If not set the write will be allowed. If set to 0 a write will only be allowed if the key doesn’t exist. If the index is non-zero the write will only be allowed if the key’s current version matches the version specified in the cas parameter. The default is -1.")
	return cmd
}

func kvPatchCmd() *cobra.Command {
	var (
		cas int32
	)

	cmd := &cobra.Command{
		Use:   "patch PATH [DATA]",
		Short: "Writes the data to the given path in the key-value store. (DATA format: bar=baz)",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)

			data, err := restflags.ParseArgsData(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}

			var c uint32
			if cas != -1 {
				c = utils.ToUint32(cas)
			}
			body := types.PostSecretRequest{
				Data: new(any),
				Options: &types.PostSecretOptions{
					Cas: &c,
				},
			}

			*(body.Data) = data

			resp := exit.OnErr2(common.Client().PatchSecretRequest(cmd.Context(), args[0], body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderSecretMetadataTable(resp.Data)
			}
		},
	}

	cmd.Flags().Int32Var(&cas, "cas", -1, "Specifies to use a Check-And-Set operation. If not set the write will be allowed. If set to 0 a write will only be allowed if the key doesn’t exist. If the index is non-zero the write will only be allowed if the key’s current version matches the version specified in the cas parameter. The default is -1.")
	return cmd
}

func kvDeleteCmd() *cobra.Command {
	var (
		versions []uint
	)

	cmd := &cobra.Command{
		Use:   "delete PATH",
		Short: "Deletes the data for the provided version and path in the key-value store.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(versions) == 0 {
				exit.OnErr(common.Client().DeleteSecretRequest(cmd.Context(), args[0]))
			} else {
				exit.OnErr(common.Client().DeleteSecretVersions(cmd.Context(), args[0], utils.ToUint32Array(versions)))
			}
		},
	}

	cmd.Flags().UintSliceVar(&versions, "versions", []uint{}, "Specifies the version numbers to delete. (Comma separated list of versions)")
	return cmd
}

func kvUndeleteCmd() *cobra.Command {
	var (
		versions []uint
	)

	cmd := &cobra.Command{
		Use:   "undelete PATH",
		Short: "Undeletes the data for the provided version and path in the key-value store.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnErr(common.Client().PostSecretUndelete(cmd.Context(), args[0], utils.ToUint32Array(versions)))
		},
	}

	cmd.Flags().UintSliceVar(&versions, "versions", []uint{}, "Specifies the version numbers to delete. (Comma separated list of versions)")
	_ = cmd.MarkFlagRequired("versions")
	return cmd
}

func kvDestroyCmd() *cobra.Command {
	var (
		versions []uint
	)

	cmd := &cobra.Command{
		Use:   "destroy PATH",
		Short: "Permanently removes the specified versions' data from the key-value store.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnErr(common.Client().PutSecretDestroy(cmd.Context(), args[0], utils.ToUint32Array(versions)))
		},
	}

	cmd.Flags().UintSliceVar(&versions, "versions", []uint{}, "Specifies the version numbers to delete. (Comma separated list of versions)")
	_ = cmd.MarkFlagRequired("versions")
	return cmd
}

func kvSubkeysCmd() *cobra.Command {
	var (
		version uint32
		depth   uint32
	)

	cmd := &cobra.Command{
		Use:   "subkeys PATH",
		Short: "Provides the subkeys within a secret entry that exists at the requested path.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var v *uint32
			if cmd.Flag("version").Changed {
				v = &version
			}

			var d *uint32
			if cmd.Flag("depth").Changed {
				d = &depth
			}

			resp := exit.OnErr2(common.Client().GetSecretSubkeys(cmd.Context(), args[0], d, v))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else if resp.Data != nil {
				renderSecretMetadataTable(resp.Data.Metadata)

				if resp.Data.Subkeys != nil {
					fmt.Println("Subkeys")
					table := tablewriter.NewWriter(os.Stdout)
					table.Header([]string{"Key", "Value"})
					kvs, ok := (*resp.Data.Subkeys).(map[string]any)
					if ok {
						for k, v := range kvs {
							exit.OnErr(table.Append([]string{k, fmt.Sprintf("%v", v)}))
						}
					}
					exit.OnErr(table.Render())
				}
			}
		},
	}

	cmd.Flags().Uint32Var(&version, "version", 0, "The version to return")
	cmd.Flags().Uint32Var(&depth, "depth", 0, "Deepest nesting level to provide in the output")
	return cmd
}

func renderSecretMetadataTable(data *types.SecretVersionMetadata) {
	if data == nil {
		return
	}
	createdAt := utils.DerefOrDefault(data.CreatedTime)
	deletionTime := utils.DerefOrDefault(data.DeletionTime)
	destroyed := utils.DerefOrDefault(data.Destroyed)

	var customMetadata string
	if data.CustomMetadata != nil {
		customMetadata = fmt.Sprintf("%v", *data.CustomMetadata)
	}
	version := "N/A"
	if data.Version != nil {
		version = fmt.Sprintf("%d", *data.Version)
	}

	fmt.Println("Metadata")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Key", "Value"})
	exit.OnErr(table.Bulk([][]string{
		{"Created at", createdAt},
		{"Custom metadata", customMetadata},
		{"Deletion time", deletionTime},
		{"Destroyed", fmt.Sprintf("%t", destroyed)},
		{"Version", version},
	}))
	exit.OnErr(table.Render())
}
