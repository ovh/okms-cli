//go:build unstable

package secrets

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

func kvMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage secrets metadata",
	}
	cmd.AddCommand(
		kvGetMetadataCommand(),
		kvPutMetadataCommand(),
		kvPatchMetadataCommand(),
		kvDeleteMetadataCommand(),
	)
	return cmd
}

func kvGetMetadataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get PATH",
		Short: "Retrieves path metadata from the KV store",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(common.Client().GetSecretsMetadata(cmd.Context(), args[0], false))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else if resp.Data != nil {
				createdAt := utils.DerefOrDefault(resp.Data.CreatedTime)
				casRequired := utils.DerefOrDefault(resp.Data.CasRequired)
				deleteVersionAfter := utils.DerefOrDefault(resp.Data.DeleteVersionAfter)
				updatedTime := utils.DerefOrDefault(resp.Data.UpdatedTime)

				var customMetadata string
				if resp.Data.CustomMetadata != nil {
					customMetadata = fmt.Sprintf("%v", *resp.Data.CustomMetadata)
				}
				currentVersion := "N/A"
				if resp.Data.CurrentVersion != nil {
					currentVersion = fmt.Sprintf("%d", *resp.Data.CurrentVersion)
				}
				maxVersions := "N/A"
				if resp.Data.MaxVersions != nil {
					maxVersions = fmt.Sprintf("%d", *resp.Data.MaxVersions)
				}
				oldestVersions := "N/A"
				if resp.Data.OldestVersion != nil {
					oldestVersions = fmt.Sprintf("%d", *resp.Data.OldestVersion)
				}

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Created at", createdAt},
					{"Custom metadata", customMetadata},
					{"Cas required", fmt.Sprintf("%t", casRequired)},
					{"Current version", currentVersion},
					{"Max. number of versions", maxVersions},
					{"Oldest version", oldestVersions},
					{"Delete version after", deleteVersionAfter},
					{"Updated time", updatedTime},
				})
				table.Render()

				if resp.Data.Versions != nil {
					// Sort the keys in the Versions map
					keys := make([]string, 0, len(*resp.Data.Versions))
					for key := range *resp.Data.Versions {
						keys = append(keys, key)
					}

					sort.Sort(sort.Reverse(sort.StringSlice(keys)))

					for _, k := range keys {
						v := (*resp.Data.Versions)[k]
						versionCreatedAt := utils.DerefOrDefault(v.CreatedTime)
						versionDeletionTime := utils.DerefOrDefault(v.DeletionTime)
						versionDestroyed := utils.DerefOrDefault(v.Destroyed)

						fmt.Printf("=== Version %s ===\n", k)
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"Key", "Value"})
						table.AppendBulk([][]string{
							{"Created at", versionCreatedAt},
							{"Deletion time", versionDeletionTime},
							{"Deletion time", fmt.Sprintf("%t", versionDestroyed)},
						})
						table.Render()
					}
				}
			}
		},
	}
}

func kvPutMetadataCommand() *cobra.Command {
	var (
		casRequired        bool
		maxVersions        int32
		deleteVersionAfter string
		customMetadata     map[string]string
	)

	cmd := &cobra.Command{
		Use:   "put PATH",
		Short: "Create a blank path in the key-value store or to update path configuration for a specified path.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var c *bool
			if cmd.Flag("cas-required").Changed {
				c = &casRequired
			}

			var d *string
			if cmd.Flag("delete-after").Changed {
				d = &deleteVersionAfter
			}

			var m *int32
			if cmd.Flag("max-versions").Changed {
				m = &maxVersions
			}

			var cm *map[string]interface{}
			if len(customMetadata) > 0 {
				tmp := make(map[string]interface{})
				cm = &tmp
				for k, v := range customMetadata {
					(*cm)[k] = v
				}
			}

			body := types.SecretUpdatableMetadata{
				CasRequired:        c,
				DeleteVersionAfter: d,
				MaxVersions:        m,
				CustomMetadata:     new(any),
			}
			*body.CustomMetadata = cm

			exit.OnErr(common.Client().PostSecretMetadata(cmd.Context(), args[0], body))
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "If true all keys will require the cas parameter to be set on all write requests.")
	cmd.Flags().Int32Var(&maxVersions, "max-versions", 0, "The number of versions to keep per key. This value applies to all keys, but a key's metadata setting can overwrite this value. Once a key has more than the configured allowed versions, the oldest version will be permanently deleted. ")
	cmd.Flags().StringVar(&deleteVersionAfter, "delete-after", "0s", "If set, specifies the length of time before a version is deleted.\nDate format, see: https://developer.hashicorp.com/vault/docs/concepts/duration-format")
	cmd.Flags().StringToStringVar(&customMetadata, "custom-metadata", map[string]string{}, "Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.\nThis can be specified multiple times to add multiple pieces of metadata.")
	return cmd
}

func kvPatchMetadataCommand() *cobra.Command {
	var (
		casRequired        bool
		maxVersions        int32
		deleteVersionAfter string
		customMetadata     map[string]string
	)

	cmd := &cobra.Command{
		Use:   "patch PATH",
		Short: "Patches path settings in the KV store",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var c *bool
			if cmd.Flag("cas-required").Changed {
				c = &casRequired
			}

			var d *string
			if cmd.Flag("delete-after").Changed {
				d = &deleteVersionAfter
			}

			var m *int32
			if cmd.Flag("max-versions").Changed {
				m = &maxVersions
			}

			var cm *map[string]interface{}
			if len(customMetadata) > 0 {
				tmp := make(map[string]interface{})
				cm = &tmp
				for k, v := range customMetadata {
					(*cm)[k] = v
				}
			}

			body := types.SecretUpdatableMetadata{
				CasRequired:        c,
				DeleteVersionAfter: d,
				MaxVersions:        m,
				CustomMetadata:     new(any),
			}
			*body.CustomMetadata = cm

			exit.OnErr(common.Client().PatchSecretMetadata(cmd.Context(), args[0], body))
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "If true all keys will require the cas parameter to be set on all write requests.")
	cmd.Flags().Int32Var(&maxVersions, "max-versions", 0, "The number of versions to keep per key. This value applies to all keys, but a key's metadata setting can overwrite this value. Once a key has more than the configured allowed versions, the oldest version will be permanently deleted. ")
	cmd.Flags().StringVar(&deleteVersionAfter, "delete-after", "0s", "If set, specifies the length of time before a version is deleted.\nDate format, see: https://developer.hashicorp.com/vault/docs/concepts/duration-format")
	cmd.Flags().StringToStringVar(&customMetadata, "custom-metadata", map[string]string{}, "Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.\nThis can be specified multiple times to add multiple pieces of metadata.")
	return cmd
}

func kvDeleteMetadataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete PATH",
		Short: "Deletes all versions and metadata for the provided path.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnErr(common.Client().DeleteSecretMetadata(cmd.Context(), args[0]))
		},
	}
}
