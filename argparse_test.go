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

func TestAll(t *testing.T) {
	parser := New()

	short_opt := ""
	parser.AddOption(Option{Name: "s", Nargs: 1, Callback: func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 1)
		short_opt = args[0]
	}})

	long_opt := ""
	parser.AddOption(Option{Name: "long", Nargs: 1, Callback: func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 1)
		long_opt = args[0]
	}})

	short_flag := false
	parser.AddOption(Option{Name: "S", Callback: func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 0)
		short_flag = true
	}})

	long_flag := false
	parser.AddOption(Option{Name: "long-flag", Nargs: 0, Callback: func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 0)
		long_flag = true
	}})

	opt_array := []string{}
	parser.AddOption(Option{Name: "two", Nargs: 2, Callback: func(ctx *Context, args ...string) {
		assertEqual(t, len(args), 2)
		opt_array = append(opt_array, args...)
	}})

	pos := ""
	parser.AddOption(StringPositional("pos", &pos))

	rest := make([]string, 0)
	parser.AddOption(StringRest("rest", &rest))

	assertError(t, false, parser.Parse("-s", "sval", "--long", "lval", "pos", "-S", "a", "--long-flag", "b", "--two", "foo", "bar", "c"))

	assertEqual(t, short_opt, "sval")
	assertEqual(t, long_opt, "lval")
	assertEqual(t, short_flag, true)
	assertEqual(t, long_flag, true)
	assertEqual(t, pos, "pos")
	assertSliceEqual(t, []string{"a", "b", "c"}, rest)
	assertSliceEqual(t, opt_array, []string{"foo", "bar"})
}

func TestMultiShort(t *testing.T) {
	parser := New()

	a := false
	parser.AddOption(Option{Name: "a", Callback: func(ctx *Context, args ...string) {
		a = true
	}})

	b := false
	parser.AddOption(Option{Name: "b", Callback: func(ctx *Context, args ...string) {
		b = true
	}})

	c := false
	parser.AddOption(Option{Name: "c", Callback: func(ctx *Context, args ...string) {
		c = true
	}})

	d := ""
	parser.AddOption(Option{Name: "d", Nargs: 1, Callback: func(ctx *Context, args ...string) {
		d = args[0]
	}})

	assertError(t, false, parser.Parse("-abcd", "val"))

	assertEqual(t, a, true)
	assertEqual(t, b, true)
	assertEqual(t, c, true)
	assertEqual(t, d, "val")

	assertError(t, true, parser.Parse("-dabc"))
}

func TestAbort(t *testing.T) {
	parser := New()

	parser.AddOption(Option{Name: "a"})

	rest := make([]string, 0)
	parser.Unparceable(func(ctx *Context, s string) {
		ctx.Abort()
		assertEqual(t, s, "-b")
		rest = append(rest, ctx.Remain()...)
	})

	assertError(t, false, parser.Parse("-a", "-b", "-c"))

	assertSliceEqual(t, []string{"-b", "-c"}, rest)

	edited := false
	parser.AddOption(Option{Name: "long", Nargs: 1})
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
	parser.AddOption(Option{Name: "a"})

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
	parser.AddOption(Option{Name: "m", Nargs: 2})
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
	parser.AddOption(Option{Name: "a", Callback: func(ctx *Context, args ...string) {
		edited = 0
	}})

	parser.AddSubParser("sub", sparser)

	sparser.AddOption(Option{Name: "b", Callback: func(ctx *Context, args ...string) {
		edited++
	}})

	sparser.AddOption(Option{Name: "a", Callback: func(ctx *Context, args ...string) {
		edited++
	}})

	assertError(t, false, parser.Parse("-a", "sub", "-b", "-a"))
	assertEqual(t, edited, 2)
	assertEqual(t, parser.SubParser, sparser)

	sparser.AddSubParser("subsub", ssparser)

	ssparser.AddOption(Option{Name: "a", Callback: func(ctx *Context, args ...string) {
		edited--
	}})

	ssparser.AddOption(Option{Name: "b", Callback: func(ctx *Context, args ...string) {
		edited--
	}})

	assertError(t, false, parser.Parse("-a", "sub", "-b", "-a", "subsub", "-a", "-b"))
	assertEqual(t, edited, 0)
	assertEqual(t, parser.SubParser, sparser)
	assertEqual(t, sparser.SubParser, ssparser)
}

func TestPositional(t *testing.T) {
	parser := New()

	pos := ""
	parser.AddOption(StringPositional("pos", &pos))
	assertError(t, false, parser.Parse("a"))
	assertEqual(t, "a", pos)

	pos2 := ""
	parser.AddOption(StringPositional("pos2", &pos2))
	assertError(t, false, parser.Parse("a", "b"))
	assertEqual(t, "a", pos)
	assertEqual(t, "b", pos2)

	assertError(t, true, parser.Parse("a", "b", "c"))

	rest := make([]string, 0)
	parser.AddOption(StringRest("rest", &rest))

	assertError(t, false, parser.Parse("a", "b", "c", "d"))
	assertEqual(t, "a", pos)
	assertEqual(t, "b", pos2)
	assertSliceEqual(t, []string{"c", "d"}, rest)

	assertError(t, true, parser.Parse("a", "b", "c", "-f", "d"))
}
