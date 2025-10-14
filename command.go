package cli

import (
	"errors"
	"reflect"
)

type command struct {
	name string
	desc string

	subcommands []*command
	options     []*option
	arguments   []*argument
}

func newCommand(name, desc string, cv reflect.Value) (*command, error) {
	if name == "" {
		return nil, errors.New("empty command name: " + cv.Type().String())
	}

	rv := reflect.Indirect(cv)
	rt := rv.Type()

	cmd := &command{name: name, desc: desc}

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)

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

		t := sf.Type
		if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct && implementsCommand(t) {
			v := rv.Field(i)
			if v.IsNil() {
				v.Set(reflect.New(t.Elem()))
			}

			c := v.Interface().(Command)
			subcommand, err := newCommand(c.Name(), c.Description(), v)
			if err != nil {
				return nil, err
			}

			if subcommand != nil {
				cmd.subcommands = append(cmd.subcommands, subcommand)
			}

			continue
		}

		short := sf.Tag.Get("short")
		if short == "-" {
			short = ""
		}

		long := sf.Tag.Get("long")
		if long == "-" {
			long = ""
		}

		if len(short) > 0 || len(long) > 0 {
			op, err := newOption(sf, rv.Field(i))
			if err != nil {
				return nil, err
			}

			if op != nil {
				cmd.options = append(cmd.options, op)
			}

			continue
		}

		arg, err := newArgument(sf, rv.Field(i))
		if err != nil {
			return nil, err
		}

		if arg != nil {
			cmd.arguments = append(cmd.arguments, arg)
		}
	}

	return cmd, nil
}

func implementsCommand(t reflect.Type) bool {
	commandType := reflect.TypeOf((*Command)(nil)).Elem()
	return t.Implements(commandType)
}
