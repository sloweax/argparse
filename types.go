package argparse

import (
	"errors"
	"fmt"
	"io"
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

func StringVar(name string, v **string) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = &args[0]
	}}
}

func StringAppend(name string, v *[]string) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = append(*v, args[0])
	}}
}

func StringPositional(name string, v *string) Option {
	return Option{Name: name, Positional: true, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = args[0]
	}}
}

func StringVarPositional(name string, v **string) Option {
	return Option{Name: name, Positional: true, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		*v = &args[0]
	}}
}

func StringRest(name string, v *[]string) Option {
	return Option{Name: name, Positional: true, Nargs: -1, Callback: func(ctx *Context, args ...string) {
		*v = append(*v, args[0])
	}}
}

func Sscanf(name string, format string, v ...any) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		if _, err := fmt.Sscanf(args[0], format, v...); err != nil {
			var rerr error
			if errors.Is(err, io.EOF) {
				rerr = fmt.Errorf("option %s %q is invalid", ctx.Option().String(), args[0])
			} else {
				rerr = fmt.Errorf("option %s %q is invalid: %s", ctx.Option().String(), args[0], err.Error())
			}
			ctx.AbortWithError(rerr)
		}
	}}
}

func Int(name string, v *int) Option {
	return Option{Name: name, Nargs: 1, Callback: func(ctx *Context, args ...string) {
		if num, err := strconv.Atoi(args[0]); err != nil {
			ctx.AbortWithError(fmt.Errorf("option %s %q requires an integer", ctx.Option().String(), args[0]))
		} else {
			*v = num
		}
	}}
}

func Func(name string, f func()) Option {
	return Option{Name: name, Callback: func(ctx *Context, args ...string) {
		f()
	}}
}
