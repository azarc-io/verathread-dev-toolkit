package main

import (
	"fmt"
	"github.com/azarc-io/verathread-dev-toolkit/cli/cmd"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"os"
)

var (
	version string

	rootCmd = &cobra.Command{
		Use:   "vdt",
		Short: "The Verathread development cli.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize an existing project that was cloned from the verathread-app-template repository.",
		Args:  cobra.RangeArgs(0, 0),
		RunE:  cmd.NewInitCmd().Cmd,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version number of the CLI",
		Args:  cobra.RangeArgs(0, 0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(version)
			return nil
		},
	}
)

func init() {
	rootCmd.Version = version
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(
		initCmd,
		versionCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("failed to execute command", "err", err)
		os.Exit(1)
	}
}
