package argparse

type Option struct {
	name     string
	nargs    int
	callback func(*Context, ...string)
}

func (o *Option) String() string {
	if len(o.name) == 1 {
		return "-" + o.name
	}
	return "--" + o.name
}
