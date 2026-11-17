package kmip

import (
	"errors"
	"fmt"
	"time"

	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func obtainLeaseCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "obtain-lease ID",
		Aliases: []string{"lease"},
		Short:   "Obtain a new lease time for an object",
		Long: `Obtain a new lease time for an object.

The lease time is the interval during which the client may use the object before a new lease must be obtained.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(kmipClient.ObtainLease(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			fmt.Println("Object ID:", resp.UniqueIdentifier)
			fmt.Println("Lease time:", resp.LeaseTime)
			if !resp.LastChangeDate.IsZero() {
				fmt.Println("Last change date:", resp.LastChangeDate.Format(time.RFC3339))
			}
		},
	}
}

func getUsageAllocationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get-usage-allocation ID",
		Aliases: []string{"usage-allocation"},
		Short:   "Obtain an allocation from the usage limits of a key",
		Long: `Obtain an allocation from the usage limits of a key.

The key must have a Usage Limits attribute. The allocated amount (--count) is expressed in the unit
of the key's usage limits, and is subtracted from its remaining usage count.`,
		Args: cobra.ExactArgs(1),
	}

	count := cmd.Flags().Int64("count", 0, "Amount to allocate, in the unit of the key's usage limits")
	_ = cmd.MarkFlagRequired("count")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if *count <= 0 {
			exit.OnErr(errors.New("--count must be greater than 0"))
		}
		resp := exit.OnErr2(kmipClient.GetUsageAllocation(args[0], *count).ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		fmt.Println("Usage allocated for object", resp.UniqueIdentifier)
	}

	return cmd
}
