package usage

import (
	"fmt"
)

type UsageError struct {
	err error
}

func (e UsageError) Error() string {
	return fmt.Errorf("usage: %w", e.err).Error()
}

func (e UsageError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (e UsageError) Unwrap() error {
	return e.err
}
