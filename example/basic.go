package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . -v -s abc
// verbose=true short=abc long= positional=

// $ go run . -vs xyz
// verbose=true short=xyz long= positional=

// $ go run . -s abc --long 123 xyz
// verbose=false short=abc long=123 positional=xyz

func main() {
	parser := argparse.New()

	verbose := false
	short := ""
	long := ""
	positional := ""

	parser.AddOption(argparse.Bool("v", &verbose))
	parser.AddOption(argparse.String("s", &short))
	parser.AddOption(argparse.String("long", &long))
	parser.AddOption(argparse.StringPositional("positional", &positional))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("verbose=%v short=%v long=%v positional=%v\n", verbose, short, long, positional)
}
