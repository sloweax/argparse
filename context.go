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
	err    error
}

func (c *Context) Abort() {
	c.abort = true
}

func (c *Context) AbortWithError(err error) {
	c.Abort()
	c.err = err
}

func (c *Context) parse() error {
	for c.Remaining() > 0 {
		if c.abort {
			break
		}

		opts, err := c.getOptions(c.args[c.index])
		if err != nil {
			if c.parser.unparceable != nil {
				c.parser.unparceable(c, c.args[c.index])
				c.index++
				continue
			}
			return err
		}

		if opts == nil {
			c.parser.SubParser = c.parser.subparsers[c.args[c.index]]
			c.index++
			return c.parser.SubParser.Parse(c.Remain()...)
		}

		c.index++
		for i, opt := range opts {
			if opt.nargs >= 1 && i != len(opts)-1 || c.Remaining() < opt.nargs {
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
			if c.err != nil {
				break
			}
			c.index += opt.nargs
		}
	}
	return c.err
}

func (c *Context) getOptions(val string) ([]*Option, error) {
	opts := make([]*Option, 0)

	if strings.HasPrefix(val, "--") && len(val) > 2 {
		optname := val[2:]
		opt, ok := c.parser.opts[optname]
		if !ok || len(opt.name) == 1 {
			return nil, fmt.Errorf("unknown option %q", val)
		}
		opts = append(opts, opt)
	} else if strings.HasPrefix(val, "-") && len(val) > 1 {
		for i := 1; i < len(val); i++ {
			optname := val[i : i+1]
			opt, ok := c.parser.opts[optname]
			if !ok {
				return nil, fmt.Errorf("unknown option %q", optname)
			}
			opts = append(opts, opt)
		}
	} else {
		if _, ok := c.parser.subparsers[val]; ok {
			return nil, nil
		}

		return nil, fmt.Errorf("could not parse option %q", val)
	}

	return opts, nil
}

func (c *Context) Remain() []string {
	return c.args[c.index:]
}

func (c *Context) Remaining() int {
	return len(c.args) - c.index
}
