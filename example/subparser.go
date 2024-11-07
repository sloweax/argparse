package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . --prefix cool- add -f file
// adding file cool-file

// $ go run . --prefix bad- del -f file
// deleting file bad-file

func main() {
	parser := argparse.New()
	add_parser := argparse.New()
	del_parser := argparse.New()

	parser.AddSubParser("add", add_parser)
	parser.AddSubParser("sub", del_parser)

	prefix := ""
	parser.AddOption("prefix", 1, func(ctx *argparse.Context, args ...string) {
		prefix = args[0]
	})

	file := ""
	file_func := func(ctx *argparse.Context, args ...string) {
		file = args[0]
	}

	add_parser.AddOption("f", 1, file_func)
	del_parser.AddOption("f", 1, file_func)

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch parser.SubParser {
	case add_parser:
		fmt.Printf("adding file %s\n", prefix+file)
	case del_parser:
		fmt.Printf("deleting file %s\n", prefix+file)
	default:
		fmt.Println("nothing to do")
	}
}
