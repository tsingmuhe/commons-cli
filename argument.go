package cli

import "reflect"

type argument struct {
	name string
}

func newArgument(sf reflect.StructField, rv reflect.Value) (*argument, error) {
	return &argument{name: sf.Name}, nil
}
