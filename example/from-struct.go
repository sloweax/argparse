package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

type MyStruct struct {
	Flag       bool   `alias:"f"` // alias accepts a comma separated list of names
	LongName   string `alias:"l"`
	Positional string `type:"positional"`
}

// $ go run . --long-name abc -f 123
// Flag=true
// LongName=abc
// Positional=123

func main() {
	ms := MyStruct{}
	parser := argparse.FromStruct(&ms)

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Flag=%v\n", ms.Flag)
	fmt.Printf("LongName=%v\n", ms.LongName)
	fmt.Printf("Positional=%v\n", ms.Positional)
}
