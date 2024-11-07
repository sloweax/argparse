package argparse

import (
	"os"
)

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

func (a *ArgParser) AddOption(opt Option) {
	a.opts[opt.Name] = &opt
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
