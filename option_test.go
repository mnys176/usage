package usage

import (
	"errors"
	"testing"
)

type optionAddArgTester struct {
	iArg string
	oErr error
}

func (oaat optionAddArgTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		if err := option.AddArg(oaat.iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := option.Args()
		if len(got) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(got))
		}
		if got[0] != oaat.iArg {
			t.Errorf("arg is %q but should be %q", got[0], oaat.iArg)
		}
	}
}

func (oaat optionAddArgTester) assertOptionArgsRepeated() func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		err1 := option.AddArg(oaat.iArg)
		err2 := option.AddArg(oaat.iArg)
		err3 := option.AddArg(oaat.iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := option.Args()
		if len(got) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(got))
		}
		if got[0] != oaat.iArg {
			t.Errorf("arg is %q but should be %q", got[0], oaat.iArg)
		}
	}
}

func (oaat optionAddArgTester) assertErrorEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		got := option.AddArg(oaat.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(got, oaat.oErr) {
			t.Errorf("got %q error but wanted %q", got, oaat.oErr)
		}
	}
}

type newOptionTester struct {
	iAliases     []string
	iDescription string
	oOption      Option
	oErr         error
}

func (not newOptionTester) assertOption() func(*testing.T) {
	return func(t *testing.T) {
		got, _ := NewOption(not.iAliases, not.iDescription)
		gotAliases := got.Aliases()
		oOptionAliases := not.oOption.Aliases()
		if len(gotAliases) != len(oOptionAliases) {
			t.Fatalf("%d aliases returned but wanted %d", len(gotAliases), len(oOptionAliases))
		}
		for i := range gotAliases {
			if gotAliases[i] != oOptionAliases[i] {
				t.Errorf("alias is %q but should be %q", gotAliases[i], oOptionAliases[i])
			}
		}
		gotDescription := got.Description()
		oOptionDescription := not.oOption.Description()
		if gotDescription != oOptionDescription {
			t.Errorf("description is %q but should be %q", gotDescription, oOptionDescription)
		}
		gotArgs := got.Args()
		oOptionArgs := not.oOption.Args()
		if len(gotArgs) != len(oOptionArgs) {
			t.Fatalf("%d args returned but wanted %d", len(gotArgs), len(oOptionArgs))
		}
		for i := range gotArgs {
			if gotArgs[i] != oOptionArgs[i] {
				t.Errorf("arg is %q but should be %q", gotArgs[i], oOptionArgs[i])
			}
		}
	}
}

func (not newOptionTester) assertErrorNoOptionAliasProvided() func(*testing.T) {
	return func(t *testing.T) {
		_, got := NewOption(not.iAliases, not.iDescription)
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(got, not.oErr) {
			t.Errorf("got %q error but wanted %q", got, not.oErr)
		}
	}
}

func (not newOptionTester) assertErrorEmptyOptionAliasString() func(*testing.T) {
	return func(t *testing.T) {
		_, got := NewOption(not.iAliases, not.iDescription)
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(got, not.oErr) {
			t.Errorf("got %q error but wanted %q", got, not.oErr)
		}
	}
}

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", optionAddArgTester{
		iArg: "foo",
	}.assertOptionArgs())
	t.Run("repeated arg strings", optionAddArgTester{
		iArg: "foo",
	}.assertOptionArgsRepeated())
	t.Run("empty arg string", optionAddArgTester{
		oErr: makeError("usage: arg string must not be empty"),
	}.assertErrorEmptyArgString())
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", newOptionTester{
		iAliases:     []string{"foo", "bar"},
		iDescription: "foo",
		oOption:      makeOption([]string{"foo", "bar"}, "foo"),
	}.assertOption())
	t.Run("single alias", newOptionTester{
		iAliases:     []string{"foo"},
		iDescription: "foo",
		oOption:      makeOption([]string{"foo"}, "foo"),
	}.assertOption())
	t.Run("single repeated alias", newOptionTester{
		iAliases:     []string{"foo", "foo"},
		iDescription: "foo",
		oOption:      makeOption([]string{"foo", "foo"}, "foo"),
	}.assertOption())
	t.Run("multiple repeated aliases", newOptionTester{
		iAliases:     []string{"foo", "bar", "foo", "bar"},
		iDescription: "foo",
		oOption:      makeOption([]string{"foo", "bar", "foo", "bar"}, "foo"),
	}.assertOption())
	t.Run("empty description string", newOptionTester{
		iAliases: []string{"foo", "bar"},
		oOption:  makeOption([]string{"foo", "bar"}, ""),
	}.assertOption())
	t.Run("nil aliases", newOptionTester{
		iDescription: "foo",
		oErr:         makeError("usage: option must have at least one alias"),
	}.assertErrorNoOptionAliasProvided())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         makeError("usage: option must have at least one alias"),
	}.assertErrorNoOptionAliasProvided())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         makeError("usage: alias string must not be empty"),
	}.assertErrorEmptyOptionAliasString())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         makeError("usage: alias string must not be empty"),
	}.assertErrorEmptyOptionAliasString())
}
