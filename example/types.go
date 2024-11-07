package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sloweax/argparse"
)

// $ go run . --pair 123 321
// v1=123 v2=321

// $ go run . --pair 123 abc
// strconv.Atoi: parsing "abc": invalid syntax
// exit status 1

func main() {
	parser := argparse.New()

	v1 := 0
	v2 := 0
	parser.AddOption(IntPair("pair", &v1, &v2))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("v1=%v v2=%v\n", v1, v2)
}

func IntPair(name string, v1 *int, v2 *int) argparse.Option {
	return argparse.Option{Name: name, Nargs: 2, Callback: func(ctx *argparse.Context, args ...string) {
		if num, err := strconv.Atoi(args[0]); err != nil {
			ctx.AbortWithError(err)
			return
		} else {
			*v1 = num
		}

		if num, err := strconv.Atoi(args[1]); err != nil {
			ctx.AbortWithError(err)
			return
		} else {
			*v2 = num
		}
	}}
}
