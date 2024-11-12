package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . --prefix cool- add file
// adding file cool-file

// $ go run . --prefix very- del bad-file
// deleting file very-bad-file

func main() {
	parser := argparse.New()
	add_parser := argparse.New()
	del_parser := argparse.New()

	parser.AddSubParser("add", add_parser)
	parser.AddSubParser("del", del_parser)

	prefix := ""
	parser.AddOption(argparse.String("prefix", &prefix))

	file := ""
	add_parser.AddOption(argparse.StringPositional("file", &file).SetRequired(true))
	del_parser.AddOption(argparse.StringPositional("file", &file).SetRequired(true))

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
