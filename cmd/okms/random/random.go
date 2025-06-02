//go:build unstable

package random

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func CreateCommand(cust common.CustomizeFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "random LENGTH",
		Aliases: []string{"rand"},
		Short:   "Generate random bytes sequence",
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}
	common.SetupRestApiFlags(cmd, cust)

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	lengthStr := args[0]
	length, err := strconv.ParseInt(lengthStr, 10, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Incorrect LENGTH argument")
		os.Exit(1)
	}

	resp := exit.OnErr2(common.Client().GenerateRandomBytes(cmd.Context(), int(length)))

	if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
		output.JsonPrint(resp)
	} else {
		r := utils.DerefOrDefault(resp.Bytes)
		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{fmt.Sprintf("Value (length: %d)", length)})
		exit.OnErr(table.Append([]string{r}))
		exit.OnErr(table.Render())
	}
}
