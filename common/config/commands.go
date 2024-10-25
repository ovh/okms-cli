package config

import (
	"errors"
	"fmt"

	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func SetProfileCommand(defaultConfig string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-profile [PROFILE]",
		Short: "Switch the default profile",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			configFile := cmd.Flag("config").Value.String()
			file, err := LoadFromFile(defaultConfig, configFile)
			if err != nil {
				exit.OnErr(fmt.Errorf("Failed to load config file: %w", err))
			}
			profile := ""
			if len(args) > 0 {
				profile = args[0]
			} else {
				profile = exit.OnErr2(pterm.DefaultInteractiveSelect.
					WithOptions(ProfileList()).
					WithDefaultOption(CurrentProfile()).
					WithOnInterruptFunc(func() {
						exit.OnErr(errors.New("Interrupted"))
					}).
					Show("Select the new profile"))
			}

			exit.OnErr(SwitchProfile(profile))
			exit.OnErr(WriteToFile(file))
		},
	}

	return cmd
}
