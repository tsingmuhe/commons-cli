package cli

import "reflect"

type option struct {
	name string
}

func newOption(sf reflect.StructField, rv reflect.Value) (*option, error) {
	return &option{name: sf.Name}, nil
}
