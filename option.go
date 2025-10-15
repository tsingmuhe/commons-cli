package cli

import "reflect"

type option struct {
	short string
	long  string
	desc  string

	required bool
}

func newOption(sf *reflect.StructField, rv reflect.Value) (*option, error) {
	short := sf.Tag.Get("short")
	long := sf.Tag.Get("long")
	if len(short) == 0 && len(long) == 0 {
		return nil, nil
	}

	desc := sf.Tag.Get("description")
	return &option{short: short, long: long, desc: desc}, nil
}
