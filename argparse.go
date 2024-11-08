package argparse

import (
	"os"
	"reflect"
	"strings"
	"unicode"
)

type ArgParser struct {
	ctx *Context

	opts map[string]*Option
	// positionals
	pos []*Option

	subparsers map[string]*ArgParser

	// selected subparser
	SubParser *ArgParser

	unparceable func(*Context, string)
}

func New() *ArgParser {
	a := new(ArgParser)
	a.opts = map[string]*Option{}
	a.subparsers = map[string]*ArgParser{}
	return a
}

func (a *ArgParser) AddOption(opt Option) {
	if opt.Positional {
		if opt.Nargs == 0 {
			panic("cant have positional with nargs == 0")
		}
		a.pos = append(a.pos, &opt)
	} else {
		if opt.Nargs < 0 {
			panic("cant have option with nargs < 0")
		}
		a.opts[opt.Name] = &opt
	}
}

func (a *ArgParser) Parse(args ...string) error {
	a.ctx = &Context{args: args, parser: a}
	return a.ctx.parse()
}

func (a *ArgParser) ParseArgs() error {
	return a.Parse(os.Args[1:]...)

}

func (a *ArgParser) Unparceable(callback func(*Context, string)) {
	a.unparceable = callback
}

func (a *ArgParser) AddSubParser(name string, p *ArgParser) {
	a.subparsers[name] = p
}

func FromStruct(s any) *ArgParser {
	p := reflect.ValueOf(s)
	v := p.Elem()
	t := v.Type()

	parser := New()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		aliases := make([]string, 0)
		aliases = append(aliases, camelCaseToDashed(ft.Name))

		if tmp, ok := ft.Tag.Lookup("alias"); ok {
			for _, alias := range strings.Split(tmp, ",") {
				aliases = append(aliases, alias)
			}
		}

		opttype := ""
		if tmp, ok := ft.Tag.Lookup("type"); ok {
			opttype = tmp
		}

		switch ft.Type.Kind() {
		case reflect.Bool:
			for _, alias := range aliases {
				switch opttype {
				case "":
					parser.AddOption(Bool(alias, (*bool)(fv.Addr().UnsafePointer())))
				default:
					panic("unsupported type")
				}
			}
		case reflect.Int:
			for _, alias := range aliases {
				switch opttype {
				case "positional":
					parser.AddOption(IntPositional(alias, (*int)(fv.Addr().UnsafePointer())))
				case "":
					parser.AddOption(Int(alias, (*int)(fv.Addr().UnsafePointer())))
				default:
					panic("unsupported type")
				}
			}
		case reflect.String:
			for _, alias := range aliases {
				switch opttype {
				case "positional":
					parser.AddOption(StringPositional(alias, (*string)(fv.Addr().UnsafePointer())))
				case "":
					parser.AddOption(String(alias, (*string)(fv.Addr().UnsafePointer())))
				default:
					panic("unsupported type")
				}
			}
		default:
			panic("unsupported type")
		}
	}

	return parser
}

func camelCaseToDashed(a string) string {
	r := strings.Builder{}
	for i, c := range a {
		if unicode.IsUpper(c) && i != 0 {
			r.WriteRune('-')
			r.WriteRune(unicode.ToLower(c))
			continue
		}
		r.WriteRune(unicode.ToLower(c))
	}
	return r.String()
}
