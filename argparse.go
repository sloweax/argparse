package argparse

type ArgParser struct {
	ctx        *Context
	opts       map[string]*Option
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

func (a *ArgParser) AddOption(name string, nargs int, callback func(*Context, ...string)) {
	opt := &Option{nargs: nargs, callback: callback, name: name}
	a.opts[name] = opt
}

func (a *ArgParser) Parse(args ...string) error {
	a.ctx = &Context{args: args, parser: a}
	return a.ctx.parse()
}

func (a *ArgParser) Unparceable(callback func(*Context, string)) {
	a.unparceable = callback
}

func (a *ArgParser) AddSubParser(name string, p *ArgParser) {
	a.subparsers[name] = p
}
