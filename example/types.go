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
	parser.AddOption(argparse.Option{Name: "num", Nargs: 1, Callback: intHandler(&num)})

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
