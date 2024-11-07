package argparse

import (
	"fmt"
	"strconv"
)

func Bool(name string, v *bool) Option {
	return Option{Name: name, Callback: func(ctx *Context, args ...string) {
		*v = true
	}}
}

func String(name string, v *string) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = args[0]
	}}
}

func StringAppend(name string, v *[]string) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = append(*v, args[0])
	}}
}

func Int(name string, v *int) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		if num, err := strconv.Atoi(args[0]); err != nil {
			ctx.AbortWithError(fmt.Errorf("%s is not an integer", args[0]))
		} else {
			*v = num
		}
	}}
}

func IntAppend(name string, v *[]int) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		if num, err := strconv.Atoi(args[0]); err != nil {
			ctx.AbortWithError(fmt.Errorf("%s is not an integer", args[0]))
		} else {
			*v = append(*v, num)
		}
	}}
}

func Func(name string, f func()) Option {
	return Option{Name: name, Callback: func(ctx *Context, args ...string) {
		f()
	}}
}
