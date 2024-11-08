package argparse

import (
	"os"
	"reflect"
	"strings"
	"unicode"
)

var (
	ArgParserKind = reflect.TypeOf(ArgParser{}).Kind()
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

func (a *ArgParser) AddOptionWithAlias(opt Option, aliases ...string) {
	a.AddOption(opt)
	for _, alias := range aliases {
		tmp := opt
		tmp.Name = alias
		a.AddOption(tmp)
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

		name := ""
		if tmp, ok := ft.Tag.Lookup("name"); ok {
			name = tmp
		} else {
			name = camelCaseToDashed(ft.Name)
		}

		aliases := make([]string, 0)
		if tmp, ok := ft.Tag.Lookup("alias"); ok {
			aliases = append(aliases, strings.Split(tmp, ",")...)
		}

		opttype := ""
		if tmp, ok := ft.Tag.Lookup("type"); ok {
			opttype = tmp
		}

		switch ft.Type.Kind() {
		case reflect.Array, reflect.Slice:
			switch fv.Type().Elem().Kind() {
			case reflect.String:
				switch opttype {
				case "positional":
					parser.AddOption(StringRest(name, (*[]string)(fv.Addr().Elem().Addr().UnsafePointer())))
				case "":
					parser.AddOptionWithAlias(StringAppend(name, (*[]string)(fv.Addr().Elem().Addr().UnsafePointer())), aliases...)
				default:
					panic("unsupported type")
				}
			default:
				panic("unsupported type")
			}
		case reflect.Pointer:
			switch ft.Type.Elem().Kind() {
			case ArgParserKind:
				parser.AddSubParser(name, (*ArgParser)(fv.UnsafePointer()))
			case reflect.String:
				switch opttype {
				case "":
					parser.AddOptionWithAlias(StringVar(name, (**string)(fv.Addr().UnsafePointer())), aliases...)
				case "positional":
					parser.AddOption(StringVarPositional(name, (**string)(fv.Addr().UnsafePointer())))
				default:
					panic("unsupported type")
				}
			}
		case reflect.Bool:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(Bool(name, (*bool)(fv.Addr().UnsafePointer())), aliases...)
			default:
				panic("unsupported type")
			}
		case reflect.Int:
			switch opttype {
			case "positional":
				parser.AddOption(IntPositional(name, (*int)(fv.Addr().UnsafePointer())))
			case "":
				parser.AddOptionWithAlias(Int(name, (*int)(fv.Addr().UnsafePointer())), aliases...)
			default:
				panic("unsupported type")
			}
		case reflect.String:
			switch opttype {
			case "positional":
				parser.AddOption(StringPositional(name, (*string)(fv.Addr().UnsafePointer())))
			case "":
				parser.AddOptionWithAlias(String(name, (*string)(fv.Addr().UnsafePointer())), aliases...)
			default:
				panic("unsupported type")
			}
		default:
			continue
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
