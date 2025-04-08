package secretsv2

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-sdk-go/types"
)

func rowFromMetadata(meta types.SecretV2Metadata) []string {
	return []string{
		fmt.Sprintf("%t", utils.DerefOrDefault(meta.CasRequired)),
		utils.DerefOrDefault(meta.CreatedAt),
		fmt.Sprintf("%d", utils.DerefOrDefault(meta.CurrentVersion)),
		utils.DerefOrDefault(meta.DeactivateVersionAfter),
		fmt.Sprintf("%d", utils.DerefOrDefault(meta.MaxVersions)),
		fmt.Sprintf("%d", utils.DerefOrDefault(meta.OldestVersion)),
		utils.DerefOrDefault(meta.UpdatedAt),
		fmt.Sprintf("%v", *meta.CustomMetadata)}
}

func renderList(secrets *types.ListSecretV2Response) {
	tableMetadata := tablewriter.NewWriter(os.Stdout)
	fmt.Println("Metadata")
	tableMetadata.SetHeader([]string{"Path", "Cas Required", "Created at", "Current Version", "Deactivate Version After", "Max Versions", "Oldest Version", "Updated at", "Custom metadata"})
	for _, secret := range *secrets.Results {
		tableMetadata.Append(append([]string{*secret.Path}, rowFromMetadata(*secret.Metadata)...))
	}
	tableMetadata.Render()
}

func renderMetadata(path string, meta types.SecretV2Metadata) {
	fmt.Printf("Metadata: %v\n", path)
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Cas Required", "Created at", "Current Version", "Deactivate Version After", "Max Versions", "Oldest Version", "Updated at", "Custom metadata"})
	table.Append(rowFromMetadata(meta))

	table.Render()
}

func renderListMetadataVersion(secrets []types.SecretV2Version) {
	fmt.Println("Metadata")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Created at", "Deactivated at", "State"})
	for _, secret := range secrets {
		table.Append([]string{fmt.Sprintf("%b", secret.Id), secret.CreatedAt, utils.DerefOrDefault(secret.DeactivatedAt), string(secret.State)})
	}
	table.Render()
}

func renderMetadataVersion(secret types.SecretV2Version) {
	// After all it's a list of size 1
	renderListMetadataVersion([]types.SecretV2Version{secret})
}

func renderDataVersion(data map[string]interface{}) {
	fmt.Println("Data")
	tableData := tablewriter.NewWriter(os.Stdout)
	tableData.SetHeader([]string{"Key", "Value"})
	for k, v := range data {
		tableData.Append([]string{k, fmt.Sprintf("%v", v)})
	}
	tableData.Render()
}
