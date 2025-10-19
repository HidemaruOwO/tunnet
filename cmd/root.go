package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewRootCmd wires the base cobra command used by the CLI.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "bridge",
		Short: "A minimal Go CLI starter template.",
		Long: `Tunnet is a lean CLI template that demonstrates how to combine
Cobra with Fang styling, animated progress, and interactive prompts.` + "\n\n" +
			"Run subcommands like `bridge progress` or `bridge prompt` to explore the building blocks.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			msg := fmt.Sprintf(
				"%s Your CLI is ready. Add your own commands when you are!",
				color.New(color.FgHiGreen, color.Bold).Sprint("[ok]"),
			)
			_, err := fmt.Fprintln(cmd.OutOrStdout(), msg)
			return err
		},
	}

	root.SilenceUsage = true
	root.AddCommand(newProgressCmd())
	root.AddCommand(newPromptCmd())

	return root
}
