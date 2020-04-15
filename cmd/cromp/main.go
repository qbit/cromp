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
		ShortUsage: "cromp",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			NewConfig(),
			NewEditor(),
			NewLoad(),
			NewLogin(),
			NewPOS(),
			NewRegister(),
			NewSimilar(),
		},
		Exec: func(context.Context, []string) error {
			entryID, err := NewDoc()
			if err != nil {
				return err
			}

			fmt.Printf("Created entry: %s\n", entryID.String())

			return nil
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
