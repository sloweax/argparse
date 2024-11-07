package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . a b --foo bar c
// foo=bar
// files=[a b c]

func main() {
	files := make([]string, 0)

	foo := ""
	parser := argparse.New()
	parser.AddOption(argparse.Option{Name: "foo", Nargs: 1, Callback: func(ctx *argparse.Context, args ...string) {
		foo = args[0]
	}})

	parser.Unparceable(func(ctx *argparse.Context, arg string) {
		files = append(files, arg)
	})

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("foo=%v\n", foo)
	fmt.Printf("files=%v\n", files)
}
