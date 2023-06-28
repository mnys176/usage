package usage

import (
	"errors"
	"testing"
)

type initTester struct {
	iName string
	oErr  error
}

func (tester initTester) assertEntry() func(*testing.T) {
	return func(t *testing.T) {
		gotErr := Init(tester.iName)
		assertNilError(t, gotErr)
		assertDefaultEntry(t, global, stringToEntry(tester.iName))
		global = nil
	}
}

func (tester initTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		got := Init(tester.iName)
		assertEmptyNameStringError(t, got, tester.oErr)
		assertNilEntry(t, global)
	}
}

type argsTester struct {
	oArgs  []string
	oPanic error
}

func (tester argsTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{args: tester.oArgs}
		got := Args()
		assertArgs(t, got, tester.oArgs)
		global = nil
	}
}

func (tester argsTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Args()
		assertNilEntry(t, global)
	}
}

func TestInit(t *testing.T) {
	t.Run("baseline", initTester{
		iName: "foo",
	}.assertEntry())
	t.Run("empty name string", initTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
}

func TestArgs(t *testing.T) {
	t.Run("baseline", argsTester{
		oArgs: []string{"foo"},
	}.assertArgs())
	t.Run("uninitialized", argsTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}
