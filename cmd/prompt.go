package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HidemaruOwO/tunnet/internal/ui"
)

func newPromptCmd() *cobra.Command {
	var accessible bool

	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Capture interactive input using charmbracelet/huh.",
		Long: `Showcases a short multi-field form powered by charmbracelet/huh.

Enable accessible mode for environments where ANSI rendering is undesirable.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := ui.RunSurvey(cmd.Context(), accessible)
			if err != nil {
				return err
			}
			msg := fmt.Sprintf(
				"Nice to meet you, %s! We'll remember that you enjoy working with %s.",
				result.Name,
				result.FavoriteLanguage,
			)
			_, err = fmt.Fprintln(cmd.OutOrStdout(), msg)
			return err
		},
	}

	cmd.Flags().BoolVar(&accessible, "accessible", false, "Render the form with accessible prompts instead of TUI output.")

	return cmd
}
