package argparse

type Option struct {
	Name       string
	Nargs      int
	Callback   func(ctx *Context, args ...string)
	Positional bool
}

func (o *Option) String() string {
	if o.Positional {
		return o.Name
	}

	if len(o.Name) == 1 {
		return "-" + o.Name
	}

	return "--" + o.Name
}
