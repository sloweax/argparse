# Examples

you can check more examples [here](https://github.com/sloweax/argparse/tree/main/example)

basic example

```go
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
```

from struct example

```go
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
	Positional string `type:"positional" required:"true"`
	// if name is not specified. it will be auto generated based on field name
	BadNameForOption string `name:"nice-name"`
	// dont add the option below
	Ignored string `ignored:"true"`
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
```

subparser example

```go
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
```
