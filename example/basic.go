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

	parser.AddOption(argparse.Bool("v", &verbose))
	parser.AddOption(argparse.String("s", &short))
	parser.AddOption(argparse.String("long", &long))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("verbose=%v short=%v long=%v\n", verbose, short, long)
}
