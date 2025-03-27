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
				renderSecretListTable(resp.Results)
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
		// TODO Add customMetadata ? How ?
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
				casRequired = utils.DerefOrDefault(resp.Metadata.CasRequired)
				createdAt := utils.DerefOrDefault(resp.Metadata.CreatedAt)
				deactivateVersionAfter := utils.DerefOrDefault(resp.Metadata.DeactivateVersionAfter)
				maxVersions := utils.DerefOrDefault(resp.Metadata.MaxVersions)

				var customMetadata string
				if resp.Metadata.CustomMetadata != nil {
					customMetadata = fmt.Sprintf("%v", *resp.Metadata.CustomMetadata)
				}

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Cas Required", fmt.Sprintf("%t", casRequired)},
					{"Created at", createdAt},
					{"Deactivate Version After", deactivateVersionAfter},
					{"Max Versions", fmt.Sprintf("%d", maxVersions)},
					{"Custom metadata", customMetadata},
					{"Path", *resp.Path}, // Path is displayed in the metadata Table, which can be confusing since it's not a metadata
					// TODO : improve path display, maybe on top of the table `Path: ...`
				})
				table.Render()
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
				casRequired := utils.DerefOrDefault(resp.Metadata.CasRequired)
				createdAt := utils.DerefOrDefault(resp.Metadata.CreatedAt)
				deactivateVersionAfter := utils.DerefOrDefault(resp.Metadata.DeactivateVersionAfter)
				maxVersions := utils.DerefOrDefault(resp.Metadata.MaxVersions)

				var customMetadata string
				if resp.Metadata.CustomMetadata != nil {
					customMetadata = fmt.Sprintf("%v", *resp.Metadata.CustomMetadata)
				}

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Cas Required", fmt.Sprintf("%t", casRequired)},
					{"Created at", createdAt},
					{"Deactivate Version After", deactivateVersionAfter},
					{"Max Versions", fmt.Sprintf("%d", maxVersions)},
					{"Custom metadata", customMetadata},
					{"Path", *resp.Path},
				})
				table.Render()

				if includeData && resp.Version.Data != nil {
					fmt.Println("Data")
					tableData := tablewriter.NewWriter(os.Stdout)
					tableData.SetHeader([]string{"Key", "Value"})
					for k, v := range *resp.Version.Data {
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

func secretPutCmd() *cobra.Command {
	var (
		casRequired            bool
		maxVersions            uint32
		deactivateVersionAfter string
		cas                    uint32
		// TODO Add customMetadata ? How ?
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
				casRequired = utils.DerefOrDefault(resp.Metadata.CasRequired)
				createdAt := utils.DerefOrDefault(resp.Metadata.CreatedAt)
				deactivateVersionAfter := utils.DerefOrDefault(resp.Metadata.DeactivateVersionAfter)
				maxVersions := utils.DerefOrDefault(resp.Metadata.MaxVersions)

				var customMetadata string
				if resp.Metadata.CustomMetadata != nil {
					customMetadata = fmt.Sprintf("%v", *resp.Metadata.CustomMetadata)
				}

				fmt.Println("Metadata")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Key", "Value"})
				table.AppendBulk([][]string{
					{"Cas Required", fmt.Sprintf("%t", casRequired)},
					{"Created at", createdAt},
					{"Deactivate Version After", deactivateVersionAfter},
					{"Max Versions", fmt.Sprintf("%d", maxVersions)},
					{"Custom metadata", customMetadata},
					{"Path", *resp.Path},
				})
				table.Render()
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

func renderSecretListTable(data *[]types.GetSecretV2Response) {
	if data == nil {
		return
	}
	for _, v := range *data {
		renderSecretTable(v)
	}
}

func renderSecretTable(data types.GetSecretV2Response) {
	casRequired := utils.DerefOrDefault(data.Metadata.CasRequired)
	createdAt := utils.DerefOrDefault(data.Metadata.CreatedAt)
	currentVersion := utils.DerefOrDefault(data.Metadata.CurrentVersion)
	deactivateVersionAfter := utils.DerefOrDefault(data.Metadata.DeactivateVersionAfter)
	maxVersions := utils.DerefOrDefault(data.Metadata.MaxVersions)
	oldestVersion := utils.DerefOrDefault(data.Metadata.OldestVersion)
	updatedAt := utils.DerefOrDefault(data.Metadata.UpdatedAt)

	var customMetadata string
	if data.Metadata.CustomMetadata != nil {
		customMetadata = fmt.Sprintf("%v", *data.Metadata.CustomMetadata)
	}

	fmt.Println("Metadata")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	table.AppendBulk([][]string{
		{"Cas Required", fmt.Sprintf("%t", casRequired)},
		{"Created at", createdAt},
		{"Current Version", fmt.Sprintf("%d", currentVersion)},
		{"Deactivate Version After", deactivateVersionAfter},
		{"Max Versions", fmt.Sprintf("%d", maxVersions)},
		{"Oldest Version", fmt.Sprintf("%d", oldestVersion)},
		{"Updated at", updatedAt},
		{"Custom metadata", customMetadata},
	})
	table.Render()

	if data.Version.Data != nil {
		fmt.Println("Data")
		tableData := tablewriter.NewWriter(os.Stdout)
		tableData.SetHeader([]string{"Key", "Value"})
		for k, v := range *data.Version.Data {
			tableData.Append([]string{k, fmt.Sprintf("%v", v)})
		}
		tableData.Render()
	}
}
