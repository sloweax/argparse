package argparse

import (
	"reflect"
	"strings"
	"unicode"
)

type Option struct {
	Name       string
	Nargs      int
	Callback   func(ctx *Context, args ...string)
	Positional bool
}

func (o *Option) String() string {
	if o.Positional {
		return o.Name
	}

	if len(o.Name) == 1 {
		return "-" + o.Name
	}

	return "--" + o.Name
}

func OptionFromStruct(name string, s any) Option {
	p := reflect.ValueOf(s)
	v := p.Elem()
	t := v.Type()
	opt := Option{}
	opt.Name = name

	fields := make([]reflect.Value, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		if len(ft.Name) == 0 || unicode.IsLower(rune(ft.Name[0])) {
			continue
		}
		if tmp, ok := ft.Tag.Lookup("ignore"); ok {
			if strings.HasPrefix(strings.ToLower(tmp), "y") {
				continue
			}
		}
		fields = append(fields, v.Field(i))
	}

	opt.Nargs = len(fields)

	opt.Callback = func(ctx *Context, args ...string) {
		for i, v := range fields {
			switch v.Interface().(type) {
			case string:
				v.SetString(args[i])
			default:
				panic("unsupported type")
			}
		}
	}

	return opt
}
