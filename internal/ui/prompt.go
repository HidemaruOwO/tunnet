package ui

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
)

// SurveyResult captures the values collected from the interactive form.
type SurveyResult struct {
	Name             string
	FavoriteLanguage string
	NewsletterOptIn  bool
}

// RunSurvey renders the interactive form and returns the populated result.
func RunSurvey(ctx context.Context, accessible bool) (*SurveyResult, error) {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}

	result := &SurveyResult{}
	form := BuildSurveyForm(result)
	form.WithAccessible(accessible)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return result, nil
}

// BuildSurveyForm constructs the form used by the prompt command.
func BuildSurveyForm(model *SurveyResult) *huh.Form {
	if model == nil {
		model = &SurveyResult{}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your name?").
				Value(&model.Name).
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return errors.New("name is required")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Favorite programming language").
				Options(
					huh.NewOption("Go", "Go"),
					huh.NewOption("Rust", "Rust"),
					huh.NewOption("Python", "Python"),
					huh.NewOption("JavaScript", "JavaScript"),
				).
				Value(&model.FavoriteLanguage).
				Validate(func(value string) error {
					if value == "" {
						return errors.New("select a language")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Subscribe to release updates?").
				Value(&model.NewsletterOptIn),
		),
	)

	form.WithShowHelp(true)
	form.WithShowErrors(true)

	return form
}
