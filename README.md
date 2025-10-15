A Go-based command line parsing tool.

### Command-line syntax

    <myapp> [<subcommand>] [OPTIONS] [--] [ARGUMENTS]

Format Explanation:

- `myapp`: Main program name
- `[<subcommand>]`: Subcommand (optional) specifying the operation to perform
- `[OPTIONS]`: Command-specific options (optional)
- `--`: Option terminator (optional) to separate options from arguments
- `[ARGUMENTS]`: Arguments (optional) providing specific data

### Conventions

#### Subcommand

- A subcommand must immediately follow the parent command.
- No options or arguments may appear between the parent command and the subcommand.

#### Options

Options are named parameters, denoted with either a hyphen and a single-letter name (e.g., `-h`) or a double hyphen and
a multiple-letter name (e.g., `--verbose`). They may or may not also include a user-specified value (e.g., `-f foo.txt`, or `--file=foo.txt`).

The order of options, generally speaking, does not affect program semantics.Subcommand-specific options must appear after the subcommand.

Boolean options take no arguments(e.g., `-f`, or `--flag`):

- Absence of `-f`/`--flag` indicates false.
- Presence of `-f`/`--flag` indicates true
- Presence of `--no-flag` indicates false. When both `-f`/`--flag` and `--no-flag` are present, the last occurrence
  takes precedence.
- Multiple short boolean options can be combined. (e.g., `-abc` equivalent to `-a -b -c`)

#### Double Dash (`--`) Separator

- Used to signify end of options
- All subsequent arguments treated as positional arguments, even if they start with `-`

#### Arguments (Positional arguments)

- The order of positional arguments is often important. (e.g., `cp foo bar` means something different from `cp bar foo`)
- The command syntax requires options to come before positional arguments.
- Positional arguments after a subcommand belong to that subcommand.

### Example commands

### Best Practices

#### Options and positional arguments

- **Prefer options to positional arguments.** It’s a bit more typing, but it makes it much clearer what is going on. It
also makes it easier to make changes to how you accept input in the future. Sometimes when using positional arguments,
it’s impossible to add new input without breaking existing behavior or creating ambiguity.

- **Have full-length versions of all options.** you don’t have to look up the meaning of flags everywhere.

- **If you’ve got two or more positional arguments for different things, you’re probably doing something wrong.** A tool should have only one core positional argument, and all other arguments must be passed via named options.

- **Display help when passed `-h` or `--help` flags.** This also applies to subcommands which might have their own help
  text.

- **Display concise help text by default.** When myapp or myapp subcommand requires arguments to function, and is run with
  no arguments, display concise help text.

- **Don’t overload -h.** Ignore any other flags and arguments that are passed.

#### Output

**Return zero exit code on success, non-zero on failure.** Exit codes are how scripts determine whether a program
succeeded or failed, so you should report this correctly. Map the non-zero exit codes to the most important failure
modes.

**Send output to stdout.** The primary output for your command should go to stdout. Anything that is machine readable
should also go to stdout—this is where piping sends things by default.

**Send messaging to stderr.** Log messages, errors, and so on should all be sent to stderr. This means that when
commands are piped together, these messages are displayed to the user and not fed into the next command.

**By default, don’t output information that’s only understandable by the creators of the software.** If a piece of
output serves only to help you (the developer) understand what your software is doing, it almost certainly shouldn’t be
displayed to normal users by default—only in verbose mode.