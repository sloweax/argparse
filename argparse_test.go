package argparse

import (
	"testing"
)

func assertError(t *testing.T, isError bool, err error) {
	if err == nil && isError {
		t.Errorf("expected error != nil")
	} else if err != nil && !isError {
		t.Errorf("expected error != nil; got %s", err.Error())
	}
}

func assertSliceEqual[T comparable](t *testing.T, expected []T, actual []T) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("expected (%+v) is not equal to actual (%+v): len(expected)=%d len(actual)=%d",
			expected, actual, len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("expected[%d] (%+v) is not equal to actual[%d] (%+v)",
				i, expected[i],
				i, actual[i])
		}
	}
}

func assertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("expected (%+v) is not equal to actual (%+v)", expected, actual)
}

func TestBasic(t *testing.T) {
	parser := New()

	short_opt := ""
	parser.AddOption("s", 1, func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 1)
		short_opt = args[0]
	})

	long_opt := ""
	parser.AddOption("long", 1, func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 1)
		long_opt = args[0]
	})

	short_flag := false
	parser.AddOption("S", 0, func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 0)
		short_flag = true
	})

	long_flag := false
	parser.AddOption("long-flag", 0, func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 0)
		long_flag = true
	})

	opt_array := []string{}
	parser.AddOption("two", 2, func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 2)
		opt_array = append(opt_array, args...)
	})

	assertError(t, false, parser.Parse("-s", "sval", "--long", "lval", "-S", "--long-flag", "--two", "foo", "bar"))

	assertEqual(t, short_opt, "sval")
	assertEqual(t, long_opt, "lval")
	assertEqual(t, short_flag, true)
	assertEqual(t, long_flag, true)
	assertSliceEqual(t, opt_array, []string{"foo", "bar"})
}

func TestMultiShort(t *testing.T) {
	parser := New()

	a := false
	parser.AddOption("a", 0, func(ctx *Context, args ...string) {
		a = true
	})

	b := false
	parser.AddOption("b", 0, func(ctx *Context, args ...string) {
		b = true
	})

	c := false
	parser.AddOption("c", 0, func(ctx *Context, args ...string) {
		c = true
	})

	d := ""
	parser.AddOption("d", 1, func(ctx *Context, args ...string) {
		d = args[0]
	})

	assertError(t, false, parser.Parse("-abcd", "val"))

	assertEqual(t, a, true)
	assertEqual(t, b, true)
	assertEqual(t, c, true)
	assertEqual(t, d, "val")

	assertError(t, true, parser.Parse("-dabc"))
}

func TestAbort(t *testing.T) {
	parser := New()

	parser.AddOption("a", 0, func(ctx *Context, args ...string) {})

	rest := make([]string, 0)
	parser.Unparceable(func(ctx *Context, s string) {
		ctx.Abort()
		assertEqual(t, s, "-b")
		rest = append(rest, ctx.Remain()...)
	})

	assertError(t, false, parser.Parse("-a", "-b", "-c"))

	assertSliceEqual(t, []string{"-b", "-c"}, rest)

	edited := false
	parser.AddOption("long", 1, func(ctx *Context, s ...string) {})
	parser.Unparceable(func(ctx *Context, s string) {
		ctx.Abort()
		assertEqual(t, s, "-long")
		edited = true
	})

	assertError(t, false, parser.Parse("-long"))

	assertEqual(t, edited, true)
}

func TestUnparceable(t *testing.T) {
	parser := New()
	parser.AddOption("a", 0, func(ctx *Context, s ...string) {})

	edited := 0
	parser.Unparceable(func(ctx *Context, s string) {
		edited++
	})

	assertError(t, false, parser.Parse("-a"))
	assertEqual(t, edited, 0)

	assertError(t, false, parser.Parse("-abababa", "asdassd", "--a"))
	assertEqual(t, edited, 3)

	parser.Unparceable(func(ctx *Context, s string) {
		parser.ctx.Abort()
		edited++
		assertEqual(t, "-ab", s)
	})

	assertError(t, false, parser.Parse("-aaaaa", "-ab", "-ba"))
	assertEqual(t, edited, 4)

	edited = 0
	parser.AddOption("m", 2, func(ctx *Context, s ...string) {})
	parser.Unparceable(func(ctx *Context, s string) {
		asserts := []string{"-am", "1"}
		assertEqual(t, asserts[edited], s)
		edited++
	})

	assertError(t, false, parser.Parse("-am", "1"))
	assertEqual(t, edited, 2)
}

func TestSubParser(t *testing.T) {
	parser := New()
	sparser := New()
	ssparser := New()

	edited := 0
	parser.AddOption("a", 0, func(ctx *Context, s ...string) {
		edited = 0
	})

	parser.AddSubParser("sub", sparser)

	sparser.AddOption("b", 0, func(ctx *Context, s ...string) {
		edited++
	})

	sparser.AddOption("a", 0, func(ctx *Context, s ...string) {
		edited++
	})

	assertError(t, false, parser.Parse("-a", "sub", "-b", "-a"))
	assertEqual(t, edited, 2)
	assertEqual(t, parser.SubParser, sparser)

	sparser.AddSubParser("subsub", ssparser)

	ssparser.AddOption("a", 0, func(ctx *Context, s ...string) {
		edited--
	})

	ssparser.AddOption("b", 0, func(ctx *Context, s ...string) {
		edited--
	})

	assertError(t, false, parser.Parse("-a", "sub", "-b", "-a", "subsub", "-a", "-b"))
	assertEqual(t, edited, 0)
	assertEqual(t, parser.SubParser, sparser)
	assertEqual(t, sparser.SubParser, ssparser)
}
