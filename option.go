package argparse

type Option struct {
	Name     string
	Nargs    int
	Callback func(ctx *Context, args ...string)
}

func (o *Option) String() string {
	if len(o.Name) == 1 {
		return "-" + o.Name
	}
	return "--" + o.Name
}
