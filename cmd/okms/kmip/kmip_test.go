package kmip

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipserver"
	"github.com/ovh/kmip-go/kmiptest"
	"github.com/ovh/kmip-go/payloads"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const testObjectId = "d2a7c0e4-5b1f-4e0c-9a7e-3f5b2c1d0e9f"

// runCommand executes the given kmip subcommand against an in-memory KMIP server routing
// the requests to exec, and returns what the command printed on stdout.
func runCommand(t *testing.T, exec *kmipserver.BatchExecutor, output string, sub *cobra.Command, args ...string) string {
	t.Helper()
	kmipClient = kmiptest.NewClientAndServer(t, exec)

	root := &cobra.Command{Use: "kmip"}
	root.PersistentFlags().String("output", output, "")
	root.AddCommand(sub)
	root.SetArgs(args)

	stdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	outC := make(chan []byte)
	go func() {
		out, _ := io.ReadAll(r)
		outC <- out
	}()

	execErr := root.ExecuteContext(context.Background())
	require.NoError(t, w.Close())
	out := <-outC
	require.NoError(t, execErr)
	return string(out)
}

func route[Req, Resp kmip.OperationPayload](exec *kmipserver.BatchExecutor, op kmip.Operation, h func(req Req) Resp) {
	exec.Route(op, kmipserver.HandleFunc(func(_ context.Context, req Req) (Resp, error) {
		return h(req), nil
	}))
}

func TestArchiveRecover(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var archived, recovered string
	route(exec, kmip.OperationArchive, func(req *payloads.ArchiveRequestPayload) *payloads.ArchiveResponsePayload {
		archived = req.UniqueIdentifier
		return &payloads.ArchiveResponsePayload{UniqueIdentifier: req.UniqueIdentifier}
	})
	route(exec, kmip.OperationRecover, func(req *payloads.RecoverRequestPayload) *payloads.RecoverResponsePayload {
		recovered = req.UniqueIdentifier
		return &payloads.RecoverResponsePayload{UniqueIdentifier: req.UniqueIdentifier}
	})

	out := runCommand(t, exec, "text", archiveCommand(), "archive", testObjectId)
	require.Equal(t, testObjectId, archived)
	require.Equal(t, "Archived object "+testObjectId+"\n", out)

	out = runCommand(t, exec, "json", recoverCommand(), "recover", testObjectId)
	require.Equal(t, testObjectId, recovered)
	var resp payloads.RecoverResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &resp))
	require.Equal(t, testObjectId, resp.UniqueIdentifier)
}

func TestObtainLease(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	lastChange := time.Date(2026, 10, 16, 8, 30, 0, 0, time.UTC)
	route(exec, kmip.OperationObtainLease, func(req *payloads.ObtainLeaseRequestPayload) *payloads.ObtainLeaseResponsePayload {
		return &payloads.ObtainLeaseResponsePayload{
			UniqueIdentifier: req.UniqueIdentifier,
			LeaseTime:        time.Hour,
			LastChangeDate:   lastChange,
		}
	})

	out := runCommand(t, exec, "text", obtainLeaseCommand(), "obtain-lease", testObjectId)
	require.Equal(t, "Object ID: "+testObjectId+"\nLease time: 1h0m0s\nLast change date: "+lastChange.Local().Format(time.RFC3339)+"\n", out)

	out = runCommand(t, exec, "json", obtainLeaseCommand(), "lease", testObjectId)
	var resp payloads.ObtainLeaseResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &resp))
	require.Equal(t, testObjectId, resp.UniqueIdentifier)
	require.Equal(t, time.Hour, resp.LeaseTime)
	require.True(t, lastChange.Equal(resp.LastChangeDate))
}

func TestGetUsageAllocation(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var count int64
	route(exec, kmip.OperationGetUsageAllocation, func(req *payloads.GetUsageAllocationRequestPayload) *payloads.GetUsageAllocationResponsePayload {
		count = req.UsageLimitsCount
		return &payloads.GetUsageAllocationResponsePayload{UniqueIdentifier: req.UniqueIdentifier}
	})

	out := runCommand(t, exec, "text", getUsageAllocationCommand(), "get-usage-allocation", testObjectId, "--count", "5")
	require.Equal(t, int64(5), count)
	require.Equal(t, "Usage allocated for object "+testObjectId+"\n", out)
}

func TestAttributesAdd(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var added kmip.Attribute
	route(exec, kmip.OperationAddAttribute, func(req *payloads.AddAttributeRequestPayload) *payloads.AddAttributeResponsePayload {
		added = req.Attribute
		return &payloads.AddAttributeResponsePayload{UniqueIdentifier: req.UniqueIdentifier, Attribute: req.Attribute}
	})

	out := runCommand(t, exec, "json", attributesCommand(), "attributes", "add", testObjectId, "Description", "my description")
	require.Equal(t, kmip.AttributeNameDescription, added.AttributeName)
	require.Equal(t, "my description", added.AttributeValue)
	require.Nil(t, added.AttributeIndex)
	var resp payloads.AddAttributeResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &resp))
	require.Equal(t, testObjectId, resp.UniqueIdentifier)

	runCommand(t, exec, "json", attributesCommand(), "attr", "add", testObjectId, "Name", "my-key", "--index", "1")
	require.Equal(t, kmip.AttributeNameName, added.AttributeName)
	require.Equal(t, kmip.Name{NameValue: "my-key", NameType: kmip.NameTypeUninterpretedTextString}, added.AttributeValue)
	require.NotNil(t, added.AttributeIndex)
	require.Equal(t, int32(1), *added.AttributeIndex)
}

func TestAttributesList(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	names := []kmip.AttributeName{kmip.AttributeNameUniqueIdentifier, kmip.AttributeNameObjectType, kmip.AttributeNameState}
	route(exec, kmip.OperationGetAttributeList, func(req *payloads.GetAttributeListRequestPayload) *payloads.GetAttributeListResponsePayload {
		return &payloads.GetAttributeListResponsePayload{UniqueIdentifier: req.UniqueIdentifier, AttributeName: names}
	})

	out := runCommand(t, exec, "text", attributesCommand(), "attributes", "list", testObjectId)
	require.Equal(t, "Unique Identifier\nObject Type\nState\n", out)

	out = runCommand(t, exec, "json", attributesCommand(), "attributes", "ls", testObjectId)
	var resp payloads.GetAttributeListResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &resp))
	require.Equal(t, names, resp.AttributeName)
}
