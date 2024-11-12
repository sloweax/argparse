package argparse

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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
	SubParser     *ArgParser
	SubParserName string

	unparceable func(*Context, string, error)
}

func New() *ArgParser {
	a := new(ArgParser)
	a.opts = map[string]*Option{}
	a.subparsers = map[string]*ArgParser{}
	return a
}

func (a *ArgParser) AddOption(opt Option) {
	if len(opt.Name) == 0 {
		panic("cant have option without name")
	}
	if strings.HasPrefix(opt.Name, "-") && len(opt.Name) != 1 && !opt.Positional {
		panic("option name cant start with -")
	}
	if opt.Positional {
		if opt.Nargs == 0 {
			panic("cant have positional with nargs == 0")
		}
		if len(a.pos) > 0 && a.pos[len(a.pos)-1].Nargs == -1 {
			panic("positional with nargs -1 must be the last")
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
	err := a.ctx.parse()
	if err != nil {
		return err
	}

	required := make([]string, 0)
	for _, opt := range a.opts {
		if opt.Required && !opt.set {
			required = append(required, opt.String())
		}
	}

	for _, opt := range a.pos {
		if opt.Required && !opt.set {
			required = append(required, opt.String())
		}
	}

	if len(required) > 0 {
		if len(required) == 1 {
			return fmt.Errorf("option %s is required", required[0])
		}
		return fmt.Errorf("the following options are required: %s", strings.Join(required, ", "))
	}

	return nil
}

func (a *ArgParser) ParseArgs() error {
	return a.Parse(os.Args[1:]...)
}

func (a *ArgParser) Unparceable(callback func(*Context, string, error)) {
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

		if !ft.IsExported() {
			continue
		}

		if tmp, ok := ft.Tag.Lookup("ignored"); ok {
			if skip, err := strconv.ParseBool(tmp); err != nil {
				panic(err)
			} else if skip {
				continue
			}
		}

		var (
			required bool
		)

		if tmp, ok := ft.Tag.Lookup("required"); ok {
			if r, err := strconv.ParseBool(tmp); err != nil {
				panic(err)
			} else if r {
				required = r
			}
		}

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

		opttype, _ := ft.Tag.Lookup("type")

		switch fv.Interface().(type) {
		case string:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(String(name, (*string)(fv.Addr().UnsafePointer())).SetRequired(required), aliases...)
			case "positional":
				parser.AddOption(StringPositional(name, (*string)(fv.Addr().UnsafePointer())).SetRequired(required))
			default:
				panic("unsupported type")
			}
		case *string:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(StringAddr(name, (**string)(fv.Addr().UnsafePointer())).SetRequired(required), aliases...)
			case "positional":
				parser.AddOption(StringAddrPositional(name, (**string)(fv.Addr().UnsafePointer())).SetRequired(required))
			default:
				panic("unsupported type")
			}
		case []string:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(StringAppend(name, (*[]string)(fv.Addr().UnsafePointer())).SetRequired(required), aliases...)
			case "positional":
				parser.AddOption(StringAppendPositional(name, (*[]string)(fv.Addr().UnsafePointer())).SetRequired(required))
			default:
				panic("unsupported type")
			}
		case bool:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(Bool(name, (*bool)(fv.Addr().UnsafePointer())).SetRequired(required), aliases...)
			default:
				panic("unsupported type")
			}
		case int:
			switch opttype {
			case "":
				parser.AddOptionWithAlias(Int(name, (*int)(fv.Addr().UnsafePointer())).SetRequired(required), aliases...)
			case "positional":
				parser.AddOption(IntPositional(name, (*int)(fv.Addr().UnsafePointer())).SetRequired(required))
			default:
				panic("unsupported type")
			}
		case ArgParser:
			parser.AddSubParser(name, (*ArgParser)(fv.Addr().UnsafePointer()))
		case *ArgParser:
			parser.AddSubParser(name, (*ArgParser)(fv.UnsafePointer()))
		default:
			if ft.Type.Kind() == reflect.Pointer {
				fv = fv.Elem()
			}
			switch fv.Type().Kind() {
			case reflect.Struct:
				switch opttype {
				case "subparser":
					parser.AddSubParser(name, FromStruct(fv.Addr().Interface()))
				default:
					panic("unsupported type")
				}
			default:
				panic("unsupported type")
			}
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
