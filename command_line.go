package cli

import (
	"errors"
	"reflect"
)

type Command interface {
	Name() string
	Description() string
	Run() int
}

type CommandLine struct {
	root *command

	version string
}

func (c *CommandLine) Run(args []string) int {
	return 0
}

func Create(c Command, version string) (*CommandLine, error) {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return nil, errors.New("command must be a pointer to a struct")
	}

	if rv.IsNil() {
		return nil, errors.New("command must not be nil")
	}

	visited := map[reflect.Type]bool{}
	cmd, err := newCommand(rv, visited)
	if err != nil {
		return nil, err
	}

	return &CommandLine{root: cmd, version: version}, nil
}
