package argparse

import (
	"fmt"
	"strings"
)

type Context struct {
	parser *ArgParser
	index  int
	args   []string
	abort  bool
}

func (c *Context) Abort() {
	c.abort = true
}

func (c *Context) parse() error {
	for c.Remaining() != 0 {
		if c.abort {
			break
		}
		opts, err := c.getOptions(c.args[c.index])
		if err != nil {
			return err
		}
		if opts == nil {
			c.parser.SubParser = c.parser.subparsers[c.args[c.index]]
			c.index++
			return c.parser.SubParser.Parse(c.Remain()...)
		}
		c.index++
		for _, opt := range opts {
			if c.Remaining() < opt.nargs {
				if c.parser.unparceable != nil {
					c.parser.unparceable(c, c.args[c.index-1])
					break
				} else {
					return fmt.Errorf("option %q requires %d arguments", opt.String(), opt.nargs)
				}
			}
			tmp := make([]string, 0, opt.nargs)
			tmp = append(tmp, c.args[c.index:c.index+opt.nargs]...)
			opt.callback(c, tmp...)
			c.index += opt.nargs
		}
	}
	return nil
}

func (c *Context) getOptions(val string) ([]*Option, error) {
	opts := make([]*Option, 0)

	if strings.HasPrefix(val, "--") {
		optname := val[2:]
		opt, ok := c.parser.opts[optname]
		if !ok || len(opt.name) == 1 {
			if c.parser.unparceable != nil {
				c.parser.unparceable(c, val)
				return []*Option{}, nil
			} else {
				return nil, fmt.Errorf("option %q is invalid", val)
			}
		}
		opts = append(opts, opt)
	} else if strings.HasPrefix(val, "-") {
		for i := 1; i < len(val); i++ {
			optname := val[i : i+1]
			opt, ok := c.parser.opts[optname]
			if !ok {
				if c.parser.unparceable != nil {
					c.parser.unparceable(c, val)
					return []*Option{}, nil
				} else {
					return nil, fmt.Errorf("option %q is invalid", val)
				}
			}
			if opt.nargs > 0 && i != len(val)-1 {
				if c.parser.unparceable != nil {
					c.parser.unparceable(c, val)
					return []*Option{}, nil
				} else {
					return nil, fmt.Errorf("option %q requires %d arguments", opt.String(), opt.nargs)
				}
			}
			opts = append(opts, opt)
		}
	} else {
		if _, ok := c.parser.subparsers[val]; ok {
			return nil, nil
		}

		if c.parser.unparceable != nil {
			c.parser.unparceable(c, val)
			return []*Option{}, nil
		} else {
			return nil, fmt.Errorf("could not parse option %q", val)
		}
	}

	return opts, nil
}

func (c *Context) Remain() []string {
	return c.args[c.index:]
}

func (c *Context) Remaining() int {
	return len(c.args) - c.index
}
