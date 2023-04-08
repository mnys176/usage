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
		optionArgs := option.Args()
		if len(optionArgs) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(optionArgs))
		}
		if optionArgs[0] != oaat.iArg {
			t.Errorf("arg is %q but should be %q", optionArgs[0], oaat.iArg)
		}
	}
}

func (oaat optionAddArgTester) assertRepeatedOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		err1 := option.AddArg(oaat.iArg)
		err2 := option.AddArg(oaat.iArg)
		err3 := option.AddArg(oaat.iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		optionArgs := option.Args()
		if len(optionArgs) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(optionArgs))
		}
		if optionArgs[0] != oaat.iArg {
			t.Errorf("arg is %q but should be %q", optionArgs[0], oaat.iArg)
		}
	}
}

func (oaat optionAddArgTester) assertEmptyArgStringError() func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		err := option.AddArg(oaat.iArg)
		if err == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(err, oaat.oErr) {
			t.Errorf("got %q error but wanted %q", err, oaat.oErr)
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
		option, _ := NewOption(not.iAliases, not.iDescription)
		optionAliases := option.Aliases()
		oOptionAliases := not.oOption.Aliases()
		if len(optionAliases) != len(oOptionAliases) {
			t.Fatalf("%d aliases returned but wanted %d", len(optionAliases), len(oOptionAliases))
		}
		for i := range optionAliases {
			if optionAliases[i] != oOptionAliases[i] {
				t.Errorf("alias is %q but should be %q", optionAliases[i], oOptionAliases[i])
			}
		}
		optionDescription := option.Description()
		oOptionDescription := not.oOption.Description()
		if optionDescription != oOptionDescription {
			t.Errorf("description is %q but should be %q", optionDescription, oOptionDescription)
		}
		optionArgs := option.Args()
		oOptionArgs := not.oOption.Args()
		if len(optionArgs) != len(oOptionArgs) {
			t.Fatalf("%d args returned but wanted %d", len(optionArgs), len(oOptionArgs))
		}
		for i := range optionArgs {
			if optionArgs[i] != oOptionArgs[i] {
				t.Errorf("arg is %q but should be %q", optionArgs[i], oOptionArgs[i])
			}
		}
	}
}

func (not newOptionTester) assertNoOptionAliasProvidedError() func(*testing.T) {
	return func(t *testing.T) {
		_, err := NewOption(not.iAliases, not.iDescription)
		if err == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(err, not.oErr) {
			t.Errorf("got %q error but wanted %q", err, not.oErr)
		}
	}
}

func (not newOptionTester) assertEmptyOptionAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		_, err := NewOption(not.iAliases, not.iDescription)
		if err == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(err, not.oErr) {
			t.Errorf("got %q error but wanted %q", err, not.oErr)
		}
	}
}

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", optionAddArgTester{
		iArg: "foo",
	}.assertOptionArgs())
	t.Run("repeated arg strings", optionAddArgTester{
		iArg: "foo",
	}.assertRepeatedOptionArgs())
	t.Run("empty arg string", optionAddArgTester{
		oErr: makeError("usage: arg string must not be empty"),
	}.assertEmptyArgStringError())
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
	}.assertNoOptionAliasProvidedError())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         makeError("usage: option must have at least one alias"),
	}.assertNoOptionAliasProvidedError())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         makeError("usage: alias string must not be empty"),
	}.assertEmptyOptionAliasStringError())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         makeError("usage: alias string must not be empty"),
	}.assertEmptyOptionAliasStringError())
}
