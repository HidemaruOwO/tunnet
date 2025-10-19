package ui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func promptInput(t *testing.T, responses ...string) io.Reader {
	t.Helper()

	pr, pw := io.Pipe()

	go func() {
		defer func() {
			if err := pw.Close(); err != nil {
				panic(fmt.Sprintf("close pipe: %v", err))
			}
		}()
		for _, line := range responses {
			if _, err := fmt.Fprintln(pw, line); err != nil {
				return
			}
		}
	}()

	return pr
}

func TestRunSurveyReturnsContextError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := RunSurvey(ctx, true)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RunSurvey() error = %v, want %v", err, context.Canceled)
	}
	if res != nil {
		t.Fatalf("RunSurvey() result = %#v, want nil", res)
	}
}

func TestBuildSurveyFormAccessibleRun(t *testing.T) {
	t.Parallel()

	model := &SurveyResult{}
	form := BuildSurveyForm(model)

	var output bytes.Buffer

	form.WithAccessible(true).
		WithInput(promptInput(t, "Hidemaru", "2", "y")).
		WithOutput(&output)

	if err := form.RunWithContext(context.Background()); err != nil {
		t.Fatalf("RunWithContext() error = %v", err)
	}

	if model.Name != "Hidemaru" {
		t.Fatalf("Name = %q, want %q", model.Name, "Hidemaru")
	}
	if model.FavoriteLanguage != "Rust" {
		t.Fatalf("FavoriteLanguage = %q, want %q", model.FavoriteLanguage, "Rust")
	}
	if !model.NewsletterOptIn {
		t.Fatalf("NewsletterOptIn = %v, want true", model.NewsletterOptIn)
	}
}

func TestBuildSurveyFormAccessibleValidation(t *testing.T) {
	t.Parallel()

	model := &SurveyResult{}
	form := BuildSurveyForm(model)

	var output bytes.Buffer

	form.WithAccessible(true).
		WithInput(promptInput(t,
			"",     // invalid name
			"Taro", // valid name
			"5",    // invalid option
			"1",    // valid option
			"",     // default confirm (false)
		)).
		WithOutput(&output)

	if err := form.RunWithContext(context.Background()); err != nil {
		t.Fatalf("RunWithContext() error = %v", err)
	}

	out := output.String()
	if !strings.Contains(out, "name is required") {
		t.Fatalf("expected validation message for name, got %q", out)
	}
	if !strings.Contains(out, "Invalid: must be a number between 1 and 4") {
		t.Fatalf("expected validation message for option range, got %q", out)
	}

	if model.Name != "Taro" {
		t.Fatalf("Name = %q, want %q", model.Name, "Taro")
	}
	if model.FavoriteLanguage != "Go" {
		t.Fatalf("FavoriteLanguage = %q, want %q", model.FavoriteLanguage, "Go")
	}
	if model.NewsletterOptIn {
		t.Fatalf("NewsletterOptIn = %v, want false", model.NewsletterOptIn)
	}
}
