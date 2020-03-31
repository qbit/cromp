package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/peterbourgon/ff/v2/ffcli"
	cromp "suah.dev/cromp/internal"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("cromp", flag.ExitOnError)

		loadFlagSet = flag.NewFlagSet("cromp load", flag.ExitOnError)

		configFlagSet = flag.NewFlagSet("cromp config", flag.ExitOnError)

		//createFlagSet = flag.NewFlagSet("cromp create", flag.ExitOnError)
	)

	root := &ffcli.Command{
		ShortUsage: "cromp <subcommand>",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			&ffcli.Command{
				Name:       "config",
				ShortUsage: "cromp config key=value",
				FlagSet:    configFlagSet,
				Exec: func(ctx context.Context, args []string) error {
					if len(args) < 1 {
						return fmt.Errorf("")
					}
					return nil
				},
			},
			&ffcli.Command{
				Name:       "load",
				ShortUsage: "cromp load <file>",
				FlagSet:    loadFlagSet,
				Exec: func(ctx context.Context, args []string) error {
					if len(args) != 1 {
						return fmt.Errorf("missing file name")
					}
					header, err := cromp.ParseFileHeader(args[0])
					if err != nil {
						return err
					}

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
					fmt.Fprintf(w, "Title:\t%s\n", header.Title)
					fmt.Fprintf(w, "Author:\t%s\n", header.Author)
					fmt.Fprintf(w, "Date:\t%s\n", header.Date)
					fmt.Fprintf(w, "UUID:\t%s\n", header.UUID)
					w.Flush()

					return nil
				},
			},
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
