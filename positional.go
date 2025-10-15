package cli

import "reflect"

type positional struct {
	name string
	desc string
}

func newPositional(sf *reflect.StructField, rv reflect.Value) (*positional, error) {
	desc := sf.Tag.Get("description")
	return &positional{name: sf.Name, desc: desc}, nil
}
