package usage

import (
	"errors"
	"fmt"
)

type UsageError struct {
	Context string
	Err     error
}

func (e UsageError) Error() string {
	return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

func (e UsageError) Is(target error) bool {
	return e.Error() == target.Error()
}

func newUsageError() (e UsageError) {
	e.Context = "usage"
	return
}

func nilEntryProvidedErr() error {
	e := newUsageError()
	e.Err = errors.New("no entry provided")
	return e
}

func nilOptionProvidedErr() error {
	e := newUsageError()
	e.Err = errors.New("no option provided")
	return e
}

func noOptionAliasProvidedErr() error {
	e := newUsageError()
	e.Err = errors.New("option must have at least one alias")
	return e
}

func emptyEntryNameStringErr() error {
	e := newUsageError()
	e.Err = errors.New("name string must not be empty")
	return e
}

func emptyOptionAliasStringErr() error {
	e := newUsageError()
	e.Err = errors.New("alias string must not be empty")
	return e
}

func existingArgsErr() error {
	e := newUsageError()
	e.Err = errors.New("cannot use subcommands with args")
	return e
}

func existingEntriesErr() error {
	e := newUsageError()
	e.Err = errors.New("cannot use args with subcommands")
	return e
}

func emptyArgStringErr() error {
	e := newUsageError()
	e.Err = errors.New("arg string must not be empty")
	return e
}
