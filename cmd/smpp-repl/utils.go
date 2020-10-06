package main

import (
	"flag"
)

func makeFlags(on func(flags *flag.FlagSet)) *flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	on(flags)
	return flags
}
