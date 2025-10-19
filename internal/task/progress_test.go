package task

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestProgressConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     ProgressConfig
		wantErr error
	}{
		{
			name:    "invalid total",
			cfg:     ProgressConfig{Total: 0, Chunk: 1, Delay: time.Millisecond},
			wantErr: ErrInvalidTotal,
		},
		{
			name:    "invalid chunk",
			cfg:     ProgressConfig{Total: 5, Chunk: 0, Delay: time.Millisecond},
			wantErr: ErrInvalidChunk,
		},
		{
			name:    "negative delay",
			cfg:     ProgressConfig{Total: 5, Chunk: 1, Delay: -time.Millisecond},
			wantErr: ErrNegativeDelay,
		},
		{
			name: "valid configuration",
			cfg:  ProgressConfig{Total: 5, Chunk: 1, Delay: time.Millisecond},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildSteps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		total int
		chunk int
		want  []int
	}{
		{
			name:  "zero total",
			total: 0,
			chunk: 1,
			want:  []int{},
		},
		{
			name:  "zero chunk",
			total: 10,
			chunk: 0,
			want:  []int{},
		},
		{
			name:  "perfect division",
			total: 6,
			chunk: 3,
			want:  []int{3, 3},
		},
		{
			name:  "remainder chunk",
			total: 7,
			chunk: 3,
			want:  []int{3, 3, 1},
		},
		{
			name:  "chunk larger than total",
			total: 4,
			chunk: 10,
			want:  []int{4},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := BuildSteps(tt.total, tt.chunk)
			if len(got) != len(tt.want) {
				t.Fatalf("BuildSteps() length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("BuildSteps()[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestRunReturnsValidationError(t *testing.T) {
	t.Parallel()

	err := Run(context.Background(), &bytes.Buffer{}, ProgressConfig{
		Total: 0,
		Chunk: 1,
		Delay: 0,
	})

	if !errors.Is(err, ErrInvalidTotal) {
		t.Fatalf("Run() error = %v, want %v", err, ErrInvalidTotal)
	}
}

func TestRunReturnsContextError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Run(ctx, &bytes.Buffer{}, ProgressConfig{
		Total: 5,
		Chunk: 2,
		Delay: 0,
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want %v", err, context.Canceled)
	}
}

func TestRunCompletesWritesNewline(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := Run(context.Background(), &buf, ProgressConfig{
		Total: 4,
		Chunk: 2,
		Delay: 0,
	})

	if err != nil {
		t.Fatalf("Run() unexpected error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\n") {
		t.Fatalf("expected output to contain newline, got %q", out)
	}
}

func TestRunAllowsNilWriter(t *testing.T) {
	t.Parallel()

	err := Run(context.Background(), nil, ProgressConfig{
		Total: 3,
		Chunk: 3,
		Delay: 0,
	})

	if err != nil {
		t.Fatalf("Run() unexpected error with nil writer: %v", err)
	}
}
