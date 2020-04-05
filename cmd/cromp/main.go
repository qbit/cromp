package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v2/ffcli"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("cromp", flag.ExitOnError)

		//createFlagSet = flag.NewFlagSet("cromp create", flag.ExitOnError)
	)

	root := &ffcli.Command{
		ShortUsage: "cromp <subcommand>",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			NewConfig(),
			NewLoad(),
			NewLogin(),
			NewPOS(),
			NewRegister(),
			NewSimilar(),
		},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
