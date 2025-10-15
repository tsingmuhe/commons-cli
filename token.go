package cli

import (
	"regexp"
	"strings"
)

type tokenType int

const (
	tokenCommand    = iota
	tokenSubCommand = iota
	tokenShort
	tokenLong
	tokenArgument
	tokenSeparator
	tokenInvalid
)

var shortOptionRegex = regexp.MustCompile(`^-[a-zA-Z0-9]+$`)
var longOptionRegex = regexp.MustCompile(`^--[a-zA-Z0-9][a-zA-Z0-9_-]*(?:=[^ ]*)?$`)

type token struct {
	typ tokenType
	val []string
}

func tokenize(args []string) []token {
	if len(args) == 0 {
		return nil
	}

	tokens := make([]token, 0, len(args))
	tokens = append(tokens, token{typ: tokenCommand, val: args[0:1]})
	args = args[1:]

	var separator bool

	for _, arg := range args {
		switch {
		case separator:
			tokens = append(tokens, token{typ: tokenArgument, val: []string{arg}})
		case arg == "--":
			tokens = append(tokens, token{typ: tokenSeparator, val: []string{arg}})
			separator = true
		case arg == "-":
			tokens = append(tokens, token{typ: tokenInvalid, val: []string{arg}})
		case strings.HasPrefix(arg, "---"):
			tokens = append(tokens, token{typ: tokenInvalid, val: []string{arg}})
		case strings.HasPrefix(arg, "--"):
			if longOptionRegex.MatchString(arg) {
				before, after, found := strings.Cut(arg, "=")
				if found {
					tokens = append(tokens, token{typ: tokenLong, val: []string{before}})
					tokens = append(tokens, token{typ: tokenArgument, val: []string{after}})
				} else {
					tokens = append(tokens, token{typ: tokenLong, val: []string{before}})
				}
			} else {
				tokens = append(tokens, token{typ: tokenInvalid, val: []string{arg}})
			}
		case strings.HasPrefix(arg, "-"):
			if shortOptionRegex.MatchString(arg) {
				tokens = append(tokens, token{typ: tokenShort, val: []string{arg[0:2]}})
				for _, c := range arg[2:] {
					tokens = append(tokens, token{typ: tokenShort, val: []string{"-" + string(c)}})
				}
			} else {
				tokens = append(tokens, token{typ: tokenInvalid, val: []string{arg}})
			}
		default:
			tokens = append(tokens, token{typ: tokenArgument, val: []string{arg}})
		}
	}

	return tokens
}

func resolveSubcommandToken(tokens []token, subcommandNames []string) token {
	subCommandToken := token{typ: tokenSubCommand}

	for i, name := range subcommandNames {
		if i >= len(tokens) {
			break
		}

		tk := tokens[i]
		if name == tk.val[0] {
			subCommandToken.val = append(subCommandToken.val, name)
		} else {
			break
		}
	}

	return subCommandToken
}
