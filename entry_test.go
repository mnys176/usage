package main

import (
	"errors"
	"testing"
)

func assertArgs(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

type entryArgsTester struct {
	oArgs []string
}

func (tester entryArgsTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{args: tester.oArgs}
		got := sampleEntry.Args()
		assertArgs(t, got, tester.oArgs)
	}
}

type entryNameTester struct {
	oName string
}

func (tester entryNameTester) assertName() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: tester.oName}
		if got := sampleEntry.Name(); got != tester.oName {
			t.Errorf("name is %q but should be %q", got, tester.oName)
		}
	}
}

type entryAddArgTester struct {
	iArg string
	oErr error
}

func (tester entryAddArgTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		args := make([]string, 0, iterations)
		sampleEntry := Entry{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleEntry.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			args = append(args, tester.iArg)
		}
		assertArgs(t, sampleEntry.args, args)
	}
}

func (tester entryAddArgTester) assertEmptyArgStringErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{args: make([]string, 0)}
		got := sampleEntry.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

func TestEntryArgs(t *testing.T) {
	t.Run("baseline", entryArgsTester{
		oArgs: []string{"foo"},
	}.assertArgs())
	t.Run("multiple args", entryArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertArgs())
	t.Run("no args", entryArgsTester{
		oArgs: make([]string, 0),
	}.assertArgs())
}

func TestEntryName(t *testing.T) {
	t.Run("baseline", entryNameTester{
		oName: "foo",
	}.assertName())
}

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", entryAddArgTester{
		iArg: "foo",
	}.assertArgs())
	t.Run("empty arg string", entryAddArgTester{
		oErr: errors.New("usage: arg string must not be empty"),
	}.assertEmptyArgStringErr())
}
