package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . -v -s abc
// verbose=true short=abc long=

// $ go run . -vs xyz
// verbose=true short=xyz long=

// $ go run . -s abc --long 123
// verbose=false short=abc long=123

func main() {
	parser := argparse.New()

	verbose := false
	short := ""
	long := ""

	parser.AddOption(argparse.Option{Name: "v", Callback: func(ctx *argparse.Context, args ...string) {
		verbose = true
	}})

	parser.AddOption(argparse.Option{Name: "s", Nargs: 1, Callback: func(ctx *argparse.Context, args ...string) {
		// args is guaranteed to have length 1
		short = args[0]
	}})

	parser.AddOption(argparse.Option{Name: "long", Nargs: 1, Callback: func(ctx *argparse.Context, args ...string) {
		long = args[0]
	}})

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("verbose=%v short=%v long=%v\n", verbose, short, long)
}
