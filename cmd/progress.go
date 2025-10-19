package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/HidemaruOwO/tunnet/internal/task"
)

func newProgressCmd() *cobra.Command {
	var (
		total = 40
		chunk = 4
		delay = 65 * time.Millisecond
	)

	cmd := &cobra.Command{
		Use:   "progress",
		Short: "Demonstrate a styled progress bar.",
		Long: `Simulates a workload while rendering a progress bar using schollz/progressbar.

The command accepts tuning flags to help you experiment with timings and granularity.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg := task.ProgressConfig{
				Total: total,
				Chunk: chunk,
				Delay: delay,
			}
			return task.Run(cmd.Context(), cmd.ErrOrStderr(), cfg)
		},
	}

	cmd.Flags().IntVar(&total, "total", total, "Total units of simulated work.")
	cmd.Flags().IntVar(&chunk, "chunk", chunk, "Amount added to the bar per iteration.")
	cmd.Flags().DurationVar(&delay, "delay", delay, "Delay between progress updates.")

	return cmd
}
