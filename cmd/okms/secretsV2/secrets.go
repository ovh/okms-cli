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
	var (
		pageSize uint32
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all secrets",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			secrets := types.ListSecretV2Response{}

			for sec, err := range common.Client().ListAllSecrets(common.GetOkmsId(), &pageSize).Iter(cmd.Context()) {
				exit.OnErr(err)
				secrets = append(secrets, *sec)
			}

			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(secrets)
			} else {
				renderList(&secrets)
			}
		},
	}

	cmd.Flags().Uint32Var(&pageSize, "page-size", 100, "Number of secrets to fetch per page (between 10 and 500)")
	return cmd
}

func secretPostCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
		customMetadata         map[string]string
	)
	cmd := &cobra.Command{
		Use:     "create PATH DATA...",
		Short:   "Create a secret. Data is in key value format, a json file can also be used by adding the prefix '@' (exp: bar=baz foo=@data.json)",
		Args:    cobra.MinimumNArgs(2),
		Example: "create foo/bar zip=zap foo=@data.json | create foo/bar @data.json",
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)
			body := types.PostSecretV2Request{
				Metadata: &types.SecretV2MetadataShort{
					CustomMetadata: utils.PtrTo(types.SecretV2CustomMetadata(customMetadata)),
				},
				Path:    args[0],
				Version: types.SecretV2VersionShort{},
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
			resp := exit.OnErr2(common.Client().PostSecretV2(cmd.Context(), common.GetOkmsId(), body))
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
	cmd.Flags().StringToStringVar(&customMetadata, "custom-metadata", map[string]string{}, "Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.\nThis can be specified multiple times to add multiple pieces of metadata.")

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
			var versionPtr *uint32
			if cmd.Flag("version").Changed {
				versionPtr = &version
			}

			resp := exit.OnErr2(common.Client().GetSecretV2(cmd.Context(), common.GetOkmsId(), args[0], versionPtr, &includeData))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				renderMetadata(utils.DerefOrDefault(resp.Path), utils.DerefOrDefault(resp.Metadata))
				renderMetadataVersion(utils.DerefOrDefault(resp.Version))
				if includeData && resp.Version.Data != nil {
					// Render metadata in addition ?
					renderDataVersion(*resp.Version.Data)
				}
			}
		},
	}

	cmd.Flags().Uint32Var(&version, "version", 0, "Secret version. If not set, the latest version will be returned.")
	cmd.Flags().BoolVar(&includeData, "include-data", false, "Include the secret data. If not set they will not be returned.")
	return cmd
}

func secretPutCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
		cas                    uint32
		customMetadata         map[string]string
	)
	cmd := &cobra.Command{
		Use:     "update PATH [DATA]...",
		Short:   "Update a secret",
		Args:    cobra.MinimumNArgs(1),
		Example: "update foo/bar zip=zap bar=@data.json | update --cas-required foo/bar @data.json",
		Run: func(cmd *cobra.Command, args []string) {
			in := io.Reader(os.Stdin)
			body := types.PutSecretV2Request{
				Metadata: &types.SecretV2MetadataShort{
					CustomMetadata: utils.PtrTo(types.SecretV2CustomMetadata(customMetadata)),
				},
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
			if data != nil {
				body.Version = &types.SecretV2VersionShort{Data: &data}
			}

			resp := exit.OnErr2(common.Client().PutSecretV2(cmd.Context(), common.GetOkmsId(), args[0], c, body))
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
	cmd.Flags().StringToStringVar(&customMetadata, "custom-metadata", map[string]string{}, "Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.\nThis can be specified multiple times to add multiple pieces of metadata.")

	return cmd
}

func secretDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete PATH",
		Short: "Delete a secret and all its versions",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnErr(common.Client().DeleteSecretV2(cmd.Context(), common.GetOkmsId(), args[0]))
			fmt.Printf("Secret %s successfully deleted\n", args[0])
		},
	}
	return cmd
}
