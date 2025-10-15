package cli

import (
	"errors"
	"reflect"
)

type command struct {
	name string
	desc string
	run  func() int

	subcommands []*command
	options     []*option
	arguments   []*argument
}

var commandInterface = reflect.TypeOf((*Command)(nil)).Elem()

func newCommand(v reflect.Value) (*command, error) {
	t := v.Type()

	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct { // command must be a pointer to a struct
		return nil, nil
	}

	if !t.Implements(commandInterface) {
		return nil, nil
	}

	if v.IsNil() {
		v.Set(reflect.New(t.Elem()))
	}

	c := v.Interface().(Command)
	name := c.Name()
	if name == "" {
		return nil, errors.New("empty command name: " + t.String())
	}

	cmd := &command{name: name, desc: c.Description(), run: c.Run}

	ev := v.Elem()
	et := ev.Type()

	for i := 0; i < et.NumField(); i++ {
		sf := et.Field(i)

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

		subcommand, err := newCommand(ev.Field(i))
		if err != nil {
			return nil, err
		}
		if subcommand != nil {
			cmd.subcommands = append(cmd.subcommands, subcommand)
			continue
		}

		op, err := newOption(sf, ev.Field(i))
		if err != nil {
			return nil, err
		}
		if op != nil {
			cmd.options = append(cmd.options, op)
			continue
		}

		arg, err := newArgument(sf, ev.Field(i))
		if err != nil {
			return nil, err
		}

		if arg != nil {
			cmd.arguments = append(cmd.arguments, arg)
		}
	}

	return cmd, nil
}
