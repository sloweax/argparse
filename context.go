package argparse

import (
	"fmt"
	"strings"
)

type Context struct {
	parser *ArgParser

	// positional index
	pindex int
	index  int

	// current option
	opt *Option

	args  []string
	abort bool
	err   error
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
			if opt.Positional {
				c.index--
			}

			nargs := opt.Nargs
			if nargs < 0 {
				nargs = 1
			}

			if nargs >= 1 && i != len(opts)-1 || c.Remaining() < nargs {
				if c.parser.unparceable != nil {
					c.parser.unparceable(c, c.args[c.index-1])
					break
				} else {
					return fmt.Errorf("option %q requires %d arguments", opt.String(), nargs)
				}
			}

			tmp := make([]string, 0, nargs)
			tmp = append(tmp, c.args[c.index:c.index+nargs]...)
			if opt.Callback != nil {
				c.opt = opt
				opt.Callback(c, tmp...)
				c.opt = nil
			}

			if c.err != nil {
				break
			}

			c.index += nargs
		}
	}
	return c.err
}

func (c *Context) getOptions(val string) ([]*Option, error) {
	opts := make([]*Option, 0)

	if strings.HasPrefix(val, "--") && len(val) > 2 {
		optname := val[2:]
		opt, ok := c.parser.opts[optname]
		if !ok || len(opt.Name) == 1 {
			return nil, fmt.Errorf("unknown option %q", "--"+optname)
		}
		opts = append(opts, opt)
	} else if strings.HasPrefix(val, "-") && len(val) > 1 {
		for i := 1; i < len(val); i++ {
			optname := val[i : i+1]
			opt, ok := c.parser.opts[optname]
			if !ok {
				return nil, fmt.Errorf("unknown option %q", "-"+optname)
			}
			opts = append(opts, opt)
		}
	} else {
		if _, ok := c.parser.subparsers[val]; ok {
			return nil, nil
		}

		if c.pindex < len(c.parser.pos) {
			opts = append(opts, c.parser.pos[c.pindex])
			c.pindex++
		} else if len(c.parser.pos) > 0 && c.parser.pos[len(c.parser.pos)-1].Nargs < 0 {
			opts = append(opts, c.parser.pos[len(c.parser.pos)-1])
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

// return current option
func (c *Context) Option() *Option {
	return c.opt
}
