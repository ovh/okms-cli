package secretsv2

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
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
		fmt.Sprintf("%v", utils.DerefOrDefault(meta.CustomMetadata))}
}

func renderList(secrets *types.ListSecretV2Response) {
	tableMetadata := tablewriter.NewWriter(os.Stdout)
	fmt.Printf("Metadata (Total count : %d):\n", utils.DerefOrDefault(secrets.TotalCount))
	tableMetadata.Header([]string{"Path", "Cas Required", "Created at", "Current Version", "Deactivate Version After", "Max Versions", "Oldest Version", "Updated at", "Custom metadata"})
	for _, secret := range *secrets.Results {
		exit.OnErr(tableMetadata.Append(append([]string{*secret.Path}, rowFromMetadata(*secret.Metadata)...)))
	}
	exit.OnErr(tableMetadata.Render())
}

func renderMetadata(path string, meta types.SecretV2Metadata) {
	fmt.Printf("Metadata: %v\n", path)
	table := tablewriter.NewWriter(os.Stdout)

	table.Header([]string{"Cas Required", "Created at", "Current Version", "Deactivate Version After", "Max Versions", "Oldest Version", "Updated at", "Custom metadata"})
	exit.OnErr(table.Append(rowFromMetadata(meta)))

	exit.OnErr(table.Render())
}

func renderListMetadataVersion(secrets []types.SecretV2Version) {
	fmt.Println("Version's specific metadata ")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Id", "Created at", "Deactivated at", "State"})
	for _, secret := range secrets {
		exit.OnErr(table.Append([]string{fmt.Sprintf("%d", secret.Id), secret.CreatedAt, utils.DerefOrDefault(secret.DeactivatedAt), string(secret.State)}))
	}
	exit.OnErr(table.Render())
}

func renderMetadataVersion(secret types.SecretV2Version) {
	// After all it's a list of size 1
	renderListMetadataVersion([]types.SecretV2Version{secret})
}

func renderDataVersion(data map[string]interface{}) {
	fmt.Println("Data")
	tableData := tablewriter.NewWriter(os.Stdout)
	tableData.Header([]string{"Key", "Value"})
	for k, v := range data {
		exit.OnErr(tableData.Append([]string{k, fmt.Sprintf("%v", v)}))
	}
	exit.OnErr(tableData.Render())
}
