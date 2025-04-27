package argparse

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	BreakLineThreshold = 80
)

type ArgParser struct {
	Name        string
	Description string

	ctx *Context

	opts map[string]*Option
	// positionals
	pos []*Option

	subparsers     map[string]*ArgParser
	subparsercount int

	// selected subparser
	SubParser     *ArgParser
	SubParserName string

	unparceable func(*Context, string, error)

	optcounter int
}

func New() *ArgParser {
	a := new(ArgParser)
	a.opts = map[string]*Option{}
	a.subparsers = map[string]*ArgParser{}
	return a
}

func NewWithDefaults() *ArgParser {
	a := New()
	a.Name = os.Args[0]
	a.AddOptionWithAlias(Option{Name: "h", Description: "shows usage and exits", Callback: func(ctx *Context, args ...string) {
		fmt.Print(a.Usage())
		os.Exit(0)
	}}, "help")
	return a
}

func (a *ArgParser) AddOption(opt Option) {
	opt.sort = a.optcounter
	a.optcounter++
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
	if opt.Positional {
		panic("positional cant have alias")
	}
	a.AddOption(opt)
	opt.basealias = opt.Name
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
		if opt.Required && !opt.set && len(opt.basealias) == 0 {
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
	p.Name = a.Name + " " + name
	p.subparsercount = a.subparsercount
	a.subparsercount += 1
	a.subparsers[name] = p
}

func (a *ArgParser) LoadStruct(s any) {
	p := reflect.ValueOf(s)
	v := p.Elem()
	t := v.Type()

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
			name     string
		)

		if tmp, ok := ft.Tag.Lookup("required"); ok {
			if r, err := strconv.ParseBool(tmp); err != nil {
				panic(err)
			} else {
				required = r
			}
		}

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
		description, _ := ft.Tag.Lookup("description")
		metavar, _ := ft.Tag.Lookup("metavar")

		switch fv.Interface().(type) {
		case string:
			switch opttype {
			case "":
				a.AddOptionWithAlias(String(name, fv.Addr().Interface().(*string)).SetAll(required, description, metavar), aliases...)
			case "positional":
				a.AddOption(StringPositional(name, fv.Addr().Interface().(*string)).SetAll(required, description, metavar))
			default:
				panic("unsupported type")
			}
		case *string:
			switch opttype {
			case "":
				a.AddOptionWithAlias(StringAddr(name, fv.Addr().Interface().(**string)).SetAll(required, description, metavar), aliases...)
			case "positional":
				a.AddOption(StringAddrPositional(name, fv.Addr().Interface().(**string)).SetAll(required, description, metavar))
			default:
				panic("unsupported type")
			}
		case []string:
			switch opttype {
			case "":
				a.AddOptionWithAlias(StringAppend(name, fv.Addr().Interface().(*[]string)).SetAll(required, description, metavar), aliases...)
			case "positional":
				a.AddOption(StringAppendPositional(name, fv.Addr().Interface().(*[]string)).SetAll(required, description, metavar))
			default:
				panic("unsupported type")
			}
		case bool:
			switch opttype {
			case "":
				a.AddOptionWithAlias(Bool(name, fv.Addr().Interface().(*bool)).SetAll(required, description, metavar), aliases...)
			default:
				panic("unsupported type")
			}
		case int:
			switch opttype {
			case "":
				a.AddOptionWithAlias(Int(name, fv.Addr().Interface().(*int)).SetAll(required, description, metavar), aliases...)
			case "positional":
				a.AddOption(IntPositional(name, fv.Addr().Interface().(*int)).SetAll(required, description, metavar))
			default:
				panic("unsupported type")
			}
		case []int:
			switch opttype {
			case "positional":
				a.AddOption(IntAppendPositional(name, fv.Addr().Interface().(*[]int)).SetAll(required, description, metavar))
			default:
				panic("unsupported type")
			}
		case uint:
			switch opttype {
			case "":
				a.AddOptionWithAlias(Uint(name, fv.Addr().Interface().(*uint)).SetAll(required, description, metavar), aliases...)
			default:
				panic("unsupported type")
			}
		case ArgParser:
			a.AddSubParser(name, fv.Addr().Interface().(*ArgParser))
		case *ArgParser:
			a.AddSubParser(name, fv.Interface().(*ArgParser))
		default:
			if ft.Type.Kind() == reflect.Pointer {
				fv = fv.Elem()
			}
			switch fv.Type().Kind() {
			case reflect.Struct:
				switch opttype {
				case "subparser":
					sub := FromStruct(fv.Addr().Interface())
					sub.Description = description
					a.AddSubParser(name, sub)
				default:
					panic("unsupported type")
				}
			default:
				panic("unsupported type")
			}
		}
	}
}

func FromStruct(s any) *ArgParser {
	parser := NewWithDefaults()
	parser.LoadStruct(s)
	return parser
}

func (a *ArgParser) Options() []*Option {
	opts := make([]*Option, 0, len(a.opts)+len(a.pos))
	for _, opt := range a.opts {
		opts = append(opts, opt)
	}
	opts = append(opts, a.pos...)
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sort < opts[j].sort
	})
	return opts
}

func (a *ArgParser) Aliases() [][]*Option {
	opts := a.Options()
	if len(opts) == 0 {
		return [][]*Option{}
	}

	aliases := make([][]*Option, 0)
	alias := make([]*Option, 0)

	for _, opt := range opts {
		if len(opt.basealias) == 0 {
			if len(alias) > 0 {
				tmp := make([]*Option, 0, len(alias))
				tmp = append(tmp, alias...)
				aliases = append(aliases, tmp)
			}
			alias = alias[:0]
		}
		alias = append(alias, opt)
	}

	if len(alias) > 0 {
		aliases = append(aliases, alias)
	}

	return aliases
}

func (a *ArgParser) String() string {
	opts := a.Options()

	strs := make([]string, 0)
	for _, opt := range opts {
		if len(opt.basealias) != 0 {
			continue
		}
		strs = append(strs, opt.string())
	}

	if len(a.subparsers) > 0 {
		strs = append(strs, "<command>")
	}
	str := fmt.Sprintf("usage: %s ", a.Name)

	if len(strs) == 0 {
		return str[:len(str)-1] + "\n"
	}

	return formatString(str, len(str), BreakLineThreshold, false, strs...)
}

func (a *ArgParser) Usage() string {
	b := &strings.Builder{}
	b.WriteString(a.String())
	if len(a.Description) > 0 {
		b.WriteRune('\n')
		b.WriteString(formatString("", 0, BreakLineThreshold, false, strings.FieldsFunc(a.Description, unicode.IsSpace)...))
	}

	strs := make([]string, 0)
	aliases := a.Aliases()
	for _, alias := range aliases {
		tmp := make([]string, 0)
		for i, opt := range alias {
			if i == len(alias)-1 && opt.Nargs > 0 && !opt.Positional {
				metavar := opt.Metavar
				if len(metavar) == 0 {
					metavar = "var"
				}
				tmp = append(tmp, fmt.Sprintf("%s %s", opt.String(), metavar))
			} else {
				tmp = append(tmp, opt.String())
			}
		}
		strs = append(strs, "    "+strings.Join(tmp, ", "))
	}

	max := 0
	for _, s := range strs {
		if len(s) > max {
			max = len(s)
		}
	}
	max += 5

	if len(strs) > 0 {
		b.WriteString("\noptions:\n")
	}

	for i, s := range strs {
		if len(aliases[i][0].Description) > 0 {
			pad := max - len(s)
			for i := 0; i < pad; i++ {
				s += " "
			}
			b.WriteString(formatString(s, len(s), BreakLineThreshold, false, strings.FieldsFunc(aliases[i][0].Description, unicode.IsSpace)...))
		} else {
			b.WriteString(s)
			b.WriteRune('\n')
		}
	}

	if len(a.subparsers) > 0 {
		b.WriteString("\ncommands:\n")
	}

	subparsers := make([]string, 0, len(a.subparsers))

	for name := range a.subparsers {
		subparsers = append(subparsers, name)
	}

	sort.Slice(subparsers, func(i, j int) bool {
		is := a.subparsers[subparsers[i]]
		js := a.subparsers[subparsers[j]]
		return is.subparsercount < js.subparsercount
	})

	for _, subname := range subparsers {
		sub := a.subparsers[subname]
		str := "    " + subname
		if len(sub.Description) > 0 {
			pad := max - len(str)
			for i := 0; i < pad; i++ {
				str += " "
			}
			b.WriteString(formatString(str, len(str), BreakLineThreshold, false, strings.FieldsFunc(sub.Description, unicode.IsSpace)...))
		} else {
			b.WriteString(str)
			b.WriteRune('\n')
		}
	}

	return b.String()
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

func formatString(linestart string, padlen, breakline int, pad bool, ss ...string) string {
	line := linestart
	b := &strings.Builder{}
	for _, s := range ss {
		line += s + " "

		if (pad && len(line)+padlen > breakline) || (!pad && len(line) > breakline) {
			if pad {
				for i := 0; i < padlen; i++ {
					b.WriteRune(' ')
				}
			}
			b.WriteString(strings.TrimRightFunc(line, unicode.IsSpace))
			b.WriteRune('\n')
			pad = true
			line = ""
		}
	}

	if len(line) > 0 {
		if pad {
			for i := 0; i < padlen; i++ {
				b.WriteRune(' ')
			}
		}
		b.WriteString(line)
		b.WriteRune('\n')
	}

	return b.String()
}
