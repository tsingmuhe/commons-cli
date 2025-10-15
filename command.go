package cli

import (
	"errors"
	"reflect"
)

var commandInterface = reflect.TypeOf((*Command)(nil)).Elem()

type command struct {
	name string
	desc string
	run  func() int

	subcommands []*command
	options     []*option
	positionals []*positional
}

func (c *command) subcommandNames() []string {
	var names []string
	for _, subcommand := range c.subcommands {
		names = append(names, subcommand.name)
	}
	return names
}

func newCommand(v reflect.Value, visited map[reflect.Type]bool) (*command, error) {
	t := v.Type()

	// command type must be a pointer to a struct
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, nil
	}

	if !t.Implements(commandInterface) {
		return nil, nil
	}

	// Check for circular dependencies
	if visited[t] {
		return nil, errors.New("circular command dependency detected for type: " + t.String())
	}
	visited[t] = true

	val := v
	if val.IsNil() {
		val = reflect.New(t.Elem())
	}

	c := val.Interface().(Command)
	if c.Name() == "" {
		return nil, errors.New("command name cannot be empty for type: " + t.String())
	}

	cmd := &command{
		name: c.Name(),
		desc: c.Description(),
		run:  c.Run,
	}

	if v.IsNil() {
		v.Set(val)
	}

	err := scanStruct(v.Elem(), func(field *reflect.StructField, value reflect.Value) error {
		subcommand, err := newCommand(value, visited)
		if err != nil {
			return err
		}

		if subcommand != nil {
			cmd.subcommands = append(cmd.subcommands, subcommand)
			return nil
		}

		opt, err := newOption(field, value)
		if err != nil {
			return err
		}

		if opt != nil {
			cmd.options = append(cmd.options, opt)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	//positional arguments and subcommands are mutually exclusive and cannot be used simultaneously in the same command.
	if len(cmd.subcommands) == 0 {
		err = scanStruct(v.Elem(), func(field *reflect.StructField, value reflect.Value) error {
			pos, err := newPositional(field, value)
			if err != nil {
				return err
			}

			if pos != nil {
				cmd.positionals = append(cmd.positionals, pos)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

type scanHandler func(*reflect.StructField, reflect.Value) error

func scanStruct(v reflect.Value, handler scanHandler) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		if !sf.IsExported() {
			// Ignore unexported fields.
			continue
		}

		if sf.Anonymous {
			t := sf.Type
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}

			if t.Kind() != reflect.Struct {
				// Ignore embedded fields of non-struct types.
				continue
			}
		}

		err := handler(&sf, v.Field(i))
		if err != nil {
			return err
		}
	}

	return nil
}
