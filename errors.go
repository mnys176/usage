package main

import (
	"fmt"
)

type UsageError struct {
	Context string
	Err     error
}

func (e UsageError) Error() string {
	return fmt.Errorf("%s: %w", e.Context, e.Err).Error()
}

func (e UsageError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (e UsageError) Unwrap() error {
	return e.Err
}
