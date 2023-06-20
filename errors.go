package main

type UsageError struct {
	Msg string
}

func (e UsageError) Error() string {
	return "usage: " + e.Msg
}

func (e UsageError) Is(target error) bool {
	return e.Error() == target.Error()
}
