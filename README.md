# Usage

An easy-to-use library for customizing the format of CLI applications usage messages. The native `flag` package is great for implementing commandline interface apps, but the default help messages leave much to be desired.

First, the default usage has no way of specifying *arguments* to the application, therefore this information is missing when the usage is triggered.

Second, and this is objectively subjective, the default usage behavior is ugly. A basic application that accepts three options will yield the following usage by default.

```
Usage of example:
  -b    boolean flag (default false)
  -n int
        integer flag (default 1234)
  -s string
        string flag (default "foo")
```

What if one wanted to change the format of the default usage? This can be achieved by overriding the `Usage` method of a `flag.FlagSet`, but the new format only applies to that instance. One would have to manually replicate the style across each `flag.FlagSet` instance; this is where the `usage` package comes in.

## Getting Started

1. Install the package.

```zsh
go get github.com/mnys176/usage
```

2. This package is intended to be used alongside the native `flag` package. Scaffold a simple CLI application as you would any `flag` application.

```go
package main

import "flag"

func main() {
	flag.Parse()
}
```

3. Initialize the `usage` package using the built-in [`init`](https://go.dev/doc/effective_go#init) function. This is where one will build out the usage most of the time.

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mnys176/usage"
)

func init() {
	// Code will panic if uninitialized!
	usage.Init("example")

	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, usage.Usage())
	}
}

func main() {
	flag.Parse()
}
```

## Adding Arguments

Add some arguments to the usage message like so.

```go
// The angle brackets are explicitly included as a stylistic
// choice. This package is unopinionated by default.
usage.AddArg("<arg1>")
usage.AddArg("<arg2>")
usage.AddArg("<arg3>")
```

**NOTE:** Arguments can only be added if there are no entries present.

The usage now looks like this.

```
Usage:
    example <arg1> <arg2> <arg3>
```

## Adding Options

Add some options to the usage message like so.

```go
// Options can also have arguments, and even multiple aliases.
option1, _ := usage.NewOption([]string{"--option1"}, "the first option")
option2, _ := usage.NewOption([]string{"--option2", "--alias"}, "the second option")
option2.AddArg("<option-arg>")
usage.AddOption(option1)
usage.AddOption(option2)
```

Now the usage looks like this.

```
Usage:
    example [options] <arg1> <arg2> <arg3>

Options:
    --option1
        the first option
    --option2, --alias <option-arg>
        the second option
```

## Adding Entries (or Subcommands)

Add some entries to the usage message like so.

```go
// Entries can have both options, arguments, and even other entries.
entryOption, _ := usage.NewOption([]string{"--entry-option"}, "entry option")
entry1, _ := usage.NewEntry("entry1", "the first entry")
entry1.AddOption(entryOption)

entry2, _ := usage.NewEntry("entry2", "the second entry")
entry2.AddArg("<entry-arg>")

usage.AddEntry(entry1)
usage.AddEntry(entry2)
```

**NOTE:** Entries can only be added if there are no arguments present.

The usage now looks like this.

```
Usage:
    example <command> [options] <args>

    To learn more about the available options for each command,
    use the --help flag like so:

    example <command> --help

Commands:
    entry1
        the first entry
    entry2 <entry-arg>
        the second entry

Options:
    --option1
        the first option
    --option2, --alias <option-arg>
        the second option
```

## Managing Subcommands

When working with the `flag` package, subcommands are implemented using multiple `flag.FlagSet` instances. In this case, the `usage.Lookup` function can be used.

```go
// child.go
var (
	childFlagSet = flag.NewFlagSet("child", flag.ContinueOnError)
	childEntry   = getSomeChildEntry()
)

func getSomeChildEntry() *usage.Entry {
	childEntry, _ := usage.NewEntry("child", "some child entry")
	return childEntry
}

func init() {
	// Override subcommand usage using lookup function.
	childFlagSet.Usage = func() {
		fmt.Fprintln(os.Stdout, usage.Lookup(childEntry.Name()))
	}
}

// main.go
func init() {
	// Code will panic if uninitialized!
	usage.Init("example")

	// Build usage in main init function.
	childEntry := getSomeChildEntry()
	usage.AddEntry(childEntry)

	// Override top-level usage.
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, usage.Usage())
	}
}

func main() {
	flag.Parse()
	childFlagSet.Parse(flag.Args()[1:])
}
```

This will return the usage of a subcommand given an entry name when the `flag.FlagSet` usage is triggered.

## Setting Templates

Don't like the default templates? The default templates for entries and options can be set to custom templates using the `usage.SetEntryTemplate` and `usage.SetOptionTemplate` functions.

```go
usage.SetEntryTemplate(
    template.Must(
        template.New("").
            Parse("{{.Name}}'s template was changed."),
    ),
)
```
