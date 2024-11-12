package argparse

type Option struct {
	Name       string
	Nargs      int
	Positional bool
	Callback   func(ctx *Context, args ...string)

	Required bool

	set bool
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

func (o Option) SetRequired(val bool) Option {
	o.Required = val
	return o
}

func (o Option) SetPositional(val bool) Option {
	o.Positional = val
	return o
}
