package main

import "regexp"

type UsageError struct {
	Context string
	Err     error
}

func (e UsageError) Error() string {
	// fmt.Errorf("%s: %w", e.Context, e.Err).Error()
	return ""
}

func (e UsageError) Is(target error) bool {
	// e.Error() == target.Error()
	return false
}

func (e UsageError) Unwrap() error {
	return nil
}

type ValidationError struct {
	UsageError
	Str     string
	Pattern *regexp.Regexp
}

type ConfigurationError struct {
	UsageError
	EntryName string
}
