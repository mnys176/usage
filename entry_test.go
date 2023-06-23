package usage

import (
	"errors"
	"testing"
)

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
		got := sampleEntry.Name()
		assertName(t, got, tester.oName)
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

type entrySetNameTester struct {
	iName string
	oErr  error
}

func (tester entrySetNameTester) assertName() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: tester.iName}
		if gotErr := sampleEntry.SetName(tester.iName); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertName(t, sampleEntry.name, tester.iName)
	}
}

func (tester entrySetNameTester) assertEmptyNameStringErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: "foo"}
		got := sampleEntry.SetName(tester.iName)
		if got == nil {
			t.Fatal("no error returned with an empty name string")
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

func TestEntrySetName(t *testing.T) {
	t.Run("baseline", entrySetNameTester{
		iName: "foo",
	}.assertName())
	t.Run("empty name string", entrySetNameTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringErr())
}
