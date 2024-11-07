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
```

subparser example

```go
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
	parser.AddSubParser("del", del_parser)

	prefix := ""
	parser.AddOption(argparse.String("prefix", &prefix))

	file := ""
	add_parser.AddOption(argparse.String("file", &file))
	del_parser.AddOption(argparse.String("file", &file))

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

multi positional example

```go
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

	parser := argparse.New()

	foo := ""
	parser.AddOption(argparse.String("foo", &foo))

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
```

custom types example

```go
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sloweax/argparse"
)

// $ go run . --pair 123 321
// v1=123 v2=321

// $ go run example/types.go --pair 123 abc
// strconv.Atoi: parsing "abc": invalid syntax
// exit status 1

func main() {
	parser := argparse.New()

	v1 := 0
	v2 := 0
	parser.AddOption(IntPair("pair", &v1, &v2))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("v1=%v v2=%v\n", v1, v2)
}

func IntPair(name string, v1 *int, v2 *int) argparse.Option {
	return argparse.Option{Name: name, Nargs: 2, Callback: func(ctx *argparse.Context, args ...string) {
		if num, err := strconv.Atoi(args[0]); err != nil {
			ctx.AbortWithError(err)
			return
		} else {
			*v1 = num
		}

		if num, err := strconv.Atoi(args[1]); err != nil {
			ctx.AbortWithError(err)
			return
		} else {
			*v2 = num
		}
	}}
}
```
