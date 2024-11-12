package argparse

type Option struct {
	Name        string
	Nargs       int
	Positional  bool
	Callback    func(ctx *Context, args ...string)
	Required    bool
	Metavar     string
	Description string

	basealias string
	set       bool
	sort      int
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

func (o *Option) string() string {
	tmp := o.String()
	if o.Nargs > 0 && !o.Positional {
		metavar := o.Metavar
		if len(metavar) == 0 {
			metavar = "var"
		}
		tmp += " " + metavar
	}
	if !o.Required {
		tmp = "[" + tmp + "]"
	}
	return tmp
}

func (o Option) SetRequired(val bool) Option {
	o.Required = val
	return o
}

func (o Option) SetPositional(val bool) Option {
	o.Positional = val
	return o
}

func (o Option) SetMetavar(val string) Option {
	o.Metavar = val
	return o
}

func (o Option) SetDescription(val string) Option {
	o.Description = val
	return o
}

func (o Option) SetAll(required bool, description, metavar string) Option {
	return o.SetRequired(required).SetMetavar(metavar).SetDescription(description)
}
