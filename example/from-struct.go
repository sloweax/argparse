package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

type MyStruct struct {
	// alias accepts a comma separated list of names
	Flag       bool   `alias:"f"`
	LongName   string `alias:"l"`
	Positional string `type:"positional"`
	// if name is not specified. it will be auto generated based on field name
	BadNameForOption string `name:"nice-name"`
	// dont add the option below
	Ignored string `ignore:"true"`
}

// go run . --long-name abc -f 123 --nice-name test
// Flag=true
// LongName=abc
// Positional=123
// BadNameForOption=test

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
	fmt.Printf("BadNameForOption=%v\n", ms.BadNameForOption)
}
