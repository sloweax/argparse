package argparse

import (
	"fmt"
	"strings"
)

type Context struct {
	parser *ArgParser

	// positional index
	pindex int

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

func (c *Context) Skip() {
	if c.Remaining() == 0 {
		panic("Skip() with nothing left")
	}
	c.args = c.args[1:]
}

func (c *Context) Peek() string {
	if c.Remaining() == 0 {
		panic("Peek() with nothing left")
	}
	return c.args[0]
}

func (c *Context) Next() string {
	if c.Remaining() == 0 {
		panic("Peek() with nothing left")
	}
	r := c.args[0]
	c.args = c.args[1:]
	return r
}

func (c *Context) NextN(n int) []string {
	if c.Remaining() < n {
		panic("NextN() out of range")
	}
	tmp := make([]string, 0, n)
	tmp = append(tmp, c.args[:n]...)
	c.args = c.args[n:]
	return tmp
}

func (c *Context) parse() error {
	for c.Remaining() > 0 {
		if c.abort {
			break
		}

		if tmp := c.expand(c.args[0]); len(tmp) > 1 {
			c.args = append(tmp, c.args[1:]...)
		}

		opt, err := c.getOption(c.Peek())
		if err != nil {
			if c.parser.unparceable != nil {
				c.parser.unparceable(c, c.Peek(), err)
				c.Skip()
				continue
			}
			return err
		}

		if opt == nil {
			c.parser.SubParserName = c.Peek()
			c.parser.SubParser = c.parser.subparsers[c.Peek()]
			c.Skip()
			return c.parser.SubParser.Parse(c.Remain()...)
		}

		if !opt.Positional {
			c.Skip()
		}

		nargs := opt.Nargs
		if nargs < 0 {
			nargs = 1
		}

		if c.Remaining() < nargs {
			var suffix string
			if nargs == 1 {
				suffix = "an argument"
			} else {
				suffix = fmt.Sprintf("%d arguments", nargs)
			}
			return fmt.Errorf("option %q requires %s", opt.String(), suffix)
		}

		if opt.Callback != nil {
			c.opt = opt
			opt.Callback(c, c.NextN(nargs)...)
			c.opt = nil
		}

		if c.err != nil {
			break
		}
	}
	return c.err
}

func (c *Context) getOption(val string) (*Option, error) {
	if strings.HasPrefix(val, "--") && len(val) > 2 {
		optname := val[2:]
		opt, ok := c.parser.opts[optname]
		if !ok || len(opt.Name) == 1 {
			return nil, fmt.Errorf("unknown option %q", val)
		}
		return opt, nil
	} else if strings.HasPrefix(val, "-") && len(val) > 1 {
		optname := val[1:]
		opt, ok := c.parser.opts[optname]
		if !ok {
			return nil, fmt.Errorf("unknown option %q", val)
		}
		return opt, nil
	}

	if _, ok := c.parser.subparsers[val]; ok {
		return nil, nil
	}

	if c.pindex < len(c.parser.pos) {
		opt := c.parser.pos[c.pindex]
		c.pindex++
		return opt, nil
	} else if len(c.parser.pos) > 0 && c.parser.pos[len(c.parser.pos)-1].Nargs < 0 {
		return c.parser.pos[len(c.parser.pos)-1], nil
	}

	if strings.HasPrefix(val, "-") {
		return nil, fmt.Errorf("unknown option %q", val)
	} else {
		return nil, fmt.Errorf("unexpected operand %q", val)
	}
}

func (c *Context) Remain() []string {
	tmp := make([]string, 0, c.Remaining())
	tmp = append(tmp, c.args...)
	return tmp
}

func (c *Context) Remaining() int {
	return len(c.args)
}

// return current option
func (c *Context) Option() *Option {
	return c.opt
}

func (c *Context) expand(val string) []string {
	r := make([]string, 0)

	if strings.HasPrefix(val, "--") && len(val) > 2 && val[2] != '=' {
		if tmp := strings.SplitN(val, "=", 2); len(tmp) > 1 {
			return tmp
		}
	} else if strings.HasPrefix(val, "-") && len(val) > 1 {
		for i := 1; i < len(val); i++ {
			optname := val[i : i+1]
			r = append(r, "-"+optname)
			opt, ok := c.parser.opts[optname]
			if ((ok && opt.Nargs > 0) || optname == "-") && i != len(val)-1 {
				r = append(r, val[i+1:])
				return r
			}
		}
	}

	if len(r) == 0 {
		r = append(r, val)
	}

	return r
}
