package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . a b -f c d
// flag=true args=[a b c d]

func main() {
	parser := argparse.New()

	flag := false
	parser.AddOption(argparse.Bool("f", &flag))

	args := make([]string, 0)
	parser.AddOption(argparse.StringAppendPositional("args", &args))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("flag=%v args=%v\n", flag, args)
}
