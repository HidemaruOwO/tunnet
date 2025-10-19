package task

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	// ErrInvalidTotal indicates that the total number of units is zero or negative.
	ErrInvalidTotal = errors.New("total must be greater than zero")
	// ErrInvalidChunk indicates that the chunk size is zero or negative.
	ErrInvalidChunk = errors.New("chunk must be greater than zero")
	// ErrNegativeDelay flags a negative delay duration.
	ErrNegativeDelay = errors.New("delay must not be negative")
)

// ProgressConfig drives the simulated progress workload.
type ProgressConfig struct {
	Total int
	Chunk int
	Delay time.Duration
}

// Validate returns an error when configuration values are outside safe bounds.
func (cfg ProgressConfig) Validate() error {
	switch {
	case cfg.Total <= 0:
		return ErrInvalidTotal
	case cfg.Chunk <= 0:
		return ErrInvalidChunk
	case cfg.Delay < 0:
		return ErrNegativeDelay
	default:
		return nil
	}
}

// BuildSteps splits a total amount of work into even chunks.
func BuildSteps(total, chunk int) []int {
	if total <= 0 || chunk <= 0 {
		return []int{}
	}

	steps := make([]int, 0, (total+chunk-1)/chunk)
	remaining := total

	for remaining > 0 {
		step := chunk
		if step > remaining {
			step = remaining
		}

		steps = append(steps, step)
		remaining -= step
	}

	return steps
}

// Run renders a progress bar to the given writer using schollz/progressbar.
func Run(ctx context.Context, w io.Writer, cfg ProgressConfig) error {
	if w == nil {
		w = io.Discard
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	steps := BuildSteps(cfg.Total, cfg.Chunk)
	bar := progressbar.NewOptions(cfg.Total,
		progressbar.OptionSetWriter(w),
		progressbar.OptionSetDescription("[green]processing[reset]"),
		progressbar.OptionEnableColorCodes(true),
	)

	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := bar.Add(step); err != nil {
			return err
		}

		if cfg.Delay == 0 {
			continue
		}

		timer := time.NewTimer(cfg.Delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	_, err := fmt.Fprintln(w)
	return err
}
