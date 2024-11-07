package main

import (
	"fmt"
	"os"

	"github.com/sloweax/argparse"
)

// $ go run . --foo bar --cmd ssh root@host --foo abc
// foo=bar
// cmd=[ssh root@host --foo abc]

func main() {
	parser := argparse.New()

	foo := ""
	parser.AddOption(argparse.String("foo", &foo))

	cmd := make([]string, 0)
	parser.AddOption(argparse.Option{Name: "cmd", Callback: func(ctx *argparse.Context, args ...string) {
		ctx.Abort()
		cmd = append(cmd, ctx.Remain()...)
	}})

	if err := parser.ParseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("foo=%v\n", foo)
	fmt.Printf("cmd=%v\n", cmd)
}
