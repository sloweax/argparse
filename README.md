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

	parser.AddOption("v", 0, func(ctx *argparse.Context, args ...string) {
		verbose = true
	})

	parser.AddOption("s", 1, func(ctx *argparse.Context, args ...string) {
		// args is guaranteed to have length 1
		short = args[0]
	})

	parser.AddOption("long", 1, func(ctx *argparse.Context, args ...string) {
		long = args[0]
	})

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

	foo := ""
	parser := argparse.New()
	parser.AddOption("foo", 1, func(ctx *argparse.Context, s ...string) {
		foo = s[0]
	})

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

types example

```go
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sloweax/argparse"
)

// $ go run . --num a
// a is not an integer
// exit status 1

// $ go run . --num 34
// num=34

func main() {
	parser := argparse.New()

	num := 0
	parser.AddOption("num", 1, intHandler(&num))

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("num=%v\n", num)
}

func intHandler(loc *int) func(*argparse.Context, ...string) {
	return func(ctx *argparse.Context, args ...string) {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			ctx.AbortWithError(fmt.Errorf("%s is not an integer", args[0]))
			return
		}
		*loc = num
	}
}
```
