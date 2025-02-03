package kmip

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"github.com/ovh/kmip-go/ttlv"
	"github.com/ovh/okms-cli/common/config"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/kmipflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CustomizeFunc func(c *cobra.Command) func(*[]kmipclient.Option)

var kmipClient *kmipclient.Client

func SetupKmipFlags(command *cobra.Command, cust CustomizeFunc) {
	debug := command.PersistentFlags().BoolP("debug", "d", false, "Activate debug mode")
	// retry := command.PersistentFlags().Uint32("retry", 4, "Maximum number of HTTP retries")
	// timeout := command.PersistentFlags().Duration("timeout", okms.DefaultHTTPClientTimeout, "Timeout duration for HTTP requests")

	f := func(*[]kmipclient.Option) {}
	if cust != nil {
		f = cust(command)
	}

	config.SetupEndpointFlags(command, "kmip", func(command *cobra.Command, cfg config.EndpointConfig) {
		middlewares := []kmipclient.Middleware{
			kmipclient.CorrelationValueMiddleware(uuid.NewString),
		}
		if *debug {
			middlewares = append(middlewares, kmipclient.DebugMiddleware(os.Stderr, ttlv.MarshalXML))
		}
		opts := []kmipclient.Option{
			kmipclient.WithTlsConfig(cfg.TlsConfig("")),
			kmipclient.WithMiddlewares(middlewares...),
		}
		f(&opts)
		kmipClient = exit.OnErr2(kmipclient.Dial(
			cfg.Endpoint,
			opts...,
		))
	})
}

func NewCommand(cust CustomizeFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kmip",
		Short: "Manage kmip objects",
	}
	SetupKmipFlags(cmd, cust)

	cmd.AddCommand(
		locateCommand(),
		createCommand(),
		attributesCommand(),
		activateCommand(),
		revokeCommand(),
		destroyCommand(),
		getCommand(),
		registerCommand(),
		rekeyCommand(),
	)

	return cmd
}

func activateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "activate ID",
		Short: "Activate an object",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(kmipClient.Activate(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			fmt.Println("Activated object", resp.UniqueIdentifier)
		},
	}
}

func revokeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke ID",
		Short: "Revoke an object",
		Args:  cobra.ExactArgs(1),
	}

	reason := kmipflags.Unspecified
	cmd.Flags().Var(&reason, "reason", "Revocation reason")
	msg := cmd.Flags().String("message", "", "Optional revocation message")
	force := cmd.Flags().Bool("force", false, "Force revoke without prompting for confirmation")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if !*force {
			if ok, _ := pterm.DefaultInteractiveConfirm.Show("Revocation cannot be undone. Continue ?"); !ok {
				exit.OnErr(errors.New("Canceled"))
			}
		}
		req := kmipClient.Revoke(args[0]).WithRevocationReasonCode(kmip.RevocationReasonCode(reason))
		if *msg != "" {
			req = req.WithRevocationMessage(*msg)
		}

		resp := exit.OnErr2(req.ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		fmt.Println("Revoked object", resp.UniqueIdentifier)
	}

	return cmd
}

func destroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy ID",
		Aliases: []string{"delete", "del"},
		Short:   "Destroy an object",
		Args:    cobra.ExactArgs(1),
	}

	force := cmd.Flags().Bool("force", false, "Force deleton without prompting for confirmation")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if !*force {
			if ok, _ := pterm.DefaultInteractiveConfirm.Show("Destroy cannot be undone. Continue ?"); !ok {
				exit.OnErr(errors.New("Canceled"))
			}
		}
		resp := exit.OnErr2(kmipClient.Destroy(args[0]).ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		fmt.Println("Destroyed object", resp.UniqueIdentifier)
	}

	return cmd
}
