package cli

import "reflect"

type argument struct {
	name string
	desc string
}

func newArgument(sf reflect.StructField, rv reflect.Value) (*argument, error) {
	desc := sf.Tag.Get("description")
	return &argument{name: sf.Name, desc: desc}, nil
}
