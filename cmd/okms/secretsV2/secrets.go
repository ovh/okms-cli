package secretsv2

import (
	"fmt"
	"io"
	"os"

	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/restflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

func secretListCmd() *cobra.Command {
	// TODO: Fix
	// It generate a 500 in the CCM currently
	var (
		page_size   uint32
		page_number uint32
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all secrets",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if page_size == 0 {
				page_size = 100
			}

			resp := exit.OnErr2(common.Client().ListSecretV2(cmd.Context(), &page_size, &page_number))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else if resp.Results != nil {
				renderList(resp)
			}
		},
	}

	cmd.Flags().Uint32Var(&page_size, "page_size", 100, "Maximum number of secrets returned in one call.")
	cmd.Flags().Uint32Var(&page_number, "page_number", 1, "Number of the page to return.")
	return cmd
}

func secretPostCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
	)
	cmd := &cobra.Command{
		Use:   "create [FLAGS] PATH [DATA]",
		Short: "Create a secret",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)
			body := types.PostSecretV2Request{
				Metadata: &types.SecretV2MetadataShort{},
				Path:     args[0],
				Version:  types.SecretV2VersionShort{},
			}

			if cmd.Flag("cas-required").Changed {
				body.Metadata.CasRequired = &casRequired
			}
			if cmd.Flag("max-versions").Changed {
				body.Metadata.MaxVersions = &maxVersions
			}
			if cmd.Flag("deactivate-version-after").Changed {
				body.Metadata.DeactivateVersionAfter = &deactivateVersionAfter
			}

			data, err := restflags.ParseArgsData(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}
			body.Version.Data = &data

			// TODO Is CAS really required for PostSecretV2 since it create the secret
			resp := exit.OnErr2(common.Client().PostSecretV2(cmd.Context(), utils.PtrTo(uint32(0)), body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadata(utils.DerefOrDefault(resp.Path), utils.DerefOrDefault(resp.Metadata))
			}
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "The cas parameter will be required for all write requests if set to true")
	cmd.Flags().Uint32Var(&maxVersions, "max-versions", 10, "The number of versions to keep (10 default)")
	cmd.Flags().StringVar(&deactivateVersionAfter, "deactivate-version-after", "", "Time duration before a version is deactivated")
	return cmd
}

func secretGetCmd() *cobra.Command {
	var (
		version     uint32
		includeData bool
	)
	cmd := &cobra.Command{
		Use:   "get PATH ",
		Short: "Retrieve a secret",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(common.Client().GetSecretV2(cmd.Context(), args[0], &version, &includeData))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadata(utils.DerefOrDefault(resp.Path), utils.DerefOrDefault(resp.Metadata))
				if includeData && resp.Version.Data != nil {
					// Render metadata in addition ?
					renderDataVersion(*resp.Version.Data)
				}
			}
		},
	}

	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version. If not set, the latest version will be returned.")
	cmd.Flags().BoolVar(&includeData, "include-data", true, "Include the secret data. If not set they will be returned.")
	return cmd
}

func secretPutCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
		cas                    uint32
	)
	cmd := &cobra.Command{
		Use:   "update [FLAGS] PATH [DATA]",
		Short: "Update a secret",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)
			body := types.PutSecretV2Request{
				Metadata: &types.SecretV2MetadataShort{},
				Version:  &types.SecretV2VersionShort{},
			}
			var c *uint32
			if cmd.Flag("cas").Changed {
				c = &cas
			}

			if cmd.Flag("cas-required").Changed {
				body.Metadata.CasRequired = &casRequired
			}
			if cmd.Flag("max-versions").Changed {
				body.Metadata.MaxVersions = &maxVersions
			}
			if cmd.Flag("deactivate-version-after").Changed {
				body.Metadata.DeactivateVersionAfter = &deactivateVersionAfter
			}

			data, err := restflags.ParseArgsData(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}
			body.Version.Data = &data

			resp := exit.OnErr2(common.Client().PutSecretV2(cmd.Context(), args[0], c, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadata(utils.DerefOrDefault(resp.Path), utils.DerefOrDefault(resp.Metadata))
			}
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "The cas parameter will be required for all write requests if set to true")
	cmd.Flags().Uint32Var(&maxVersions, "max-versions", 10, "The number of versions to keep (10 default)")
	cmd.Flags().StringVar(&deactivateVersionAfter, "deactivate-version-after", "", "Time duration before a version is deactivated")
	cmd.Flags().Uint32Var(&cas, "cas", 0, "Secret version number. Required if cas-required is set to true.")

	return cmd
}

func secretPutCustomMetadataCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
		cas                    uint32
	)
	cmd := &cobra.Command{
		Use:   "update-metadata [FLAGS] PATH [CUSTOM-DATA]",
		Short: "Update a secret",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)
			body := types.PutSecretV2Request{
				Metadata: &types.SecretV2MetadataShort{},
				Version:  &types.SecretV2VersionShort{},
			}
			var c *uint32
			if cmd.Flag("cas").Changed {
				c = &cas
			}

			if cmd.Flag("cas-required").Changed {
				body.Metadata.CasRequired = &casRequired
			}
			if cmd.Flag("max-versions").Changed {
				body.Metadata.MaxVersions = &maxVersions
			}
			if cmd.Flag("deactivate-version-after").Changed {
				body.Metadata.DeactivateVersionAfter = &deactivateVersionAfter
			}

			customMetadata, err := restflags.ParseArgsCustomMetadata(in, args[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse K=V data:", err)
				os.Exit(1)
			}
			body.Metadata.CustomMetadata = utils.PtrTo(types.SecretV2CustomMetadata(customMetadata))

			resp := exit.OnErr2(common.Client().PutSecretV2(cmd.Context(), args[0], c, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadata(utils.DerefOrDefault(resp.Path), utils.DerefOrDefault(resp.Metadata))
			}
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "The cas parameter will be required for all write requests if set to true")
	cmd.Flags().Uint32Var(&maxVersions, "max-versions", 10, "The number of versions to keep (10 default)")
	cmd.Flags().StringVar(&deactivateVersionAfter, "deactivate-version-after", "", "Time duration before a version is deactivated")
	cmd.Flags().Uint32Var(&cas, "cas", 0, "Secret version number. Required if cas-required is set to true.")

	return cmd
}

func secretDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete PATH",
		Short: "Delete a secret and all its versions",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnErr(common.Client().DeleteSecretV2(cmd.Context(), args[0]))
			fmt.Printf("Secret %s successfully deleted\n", args[0])
		},
	}
	return cmd
}
