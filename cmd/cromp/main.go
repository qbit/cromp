package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
	cromp "suah.dev/cromp/internal"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("cromp", flag.ExitOnError)

		loadFlagSet = flag.NewFlagSet("cromp load", flag.ExitOnError)

		configFlagSet = flag.NewFlagSet("cromp config", flag.ExitOnError)
		urlConfFS     = configFlagSet.String("url", "", "URL of cromp server")
		tokenConfFS   = configFlagSet.String("token", "", "Access token for cromp server")

		//createFlagSet = flag.NewFlagSet("cromp create", flag.ExitOnError)
	)

	root := &ffcli.Command{
		ShortUsage: "cromp <subcommand>",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			&ffcli.Command{
				Name:       "config",
				ShortUsage: "cromp config -token [token] -url [url]",
				FlagSet:    configFlagSet,
				Exec: func(ctx context.Context, args []string) error {
					cfg := &Config{
						Token: *tokenConfFS,
						URL:   *urlConfFS,
					}

					err := cfg.WriteConfig()
					if err != nil {
						return err
					}

					return nil
				},
			},
			&ffcli.Command{
				Name:       "similar",
				ShortUsage: "cromp similar [text]",
				Exec: func(ctx context.Context, args []string) error {
					var params db.SimilarEntriesParams
					resp := &[]db.SimilarEntriesRow{}

					params.Similarity = strings.Join(args, " ")

					err := Post("/entries/similar", params, resp)
					if err != nil {
						return err
					}

					for _, e := range *resp {
						fmt.Printf("%s\t%s\t%f\n",
							e.EntryID,
							e.Title,
							e.Similarity)
					}

					return nil
				},
			},
			&ffcli.Command{
				Name:       "list",
				ShortUsage: "cromp list",
				FlagSet:    configFlagSet,
				Exec: func(ctx context.Context, args []string) error {
					resp := &[]db.Entry{}

					err := Get("/entries/list", resp)
					if err != nil {
						return err
					}

					for _, e := range *resp {
						fmt.Printf("%s\t%s\t%s\n", e.EntryID.String(),
							e.CreatedAt,
							e.Title)
					}

					return nil
				},
			},
			&ffcli.Command{
				Name:       "load",
				ShortUsage: "cromp load <file>",
				FlagSet:    loadFlagSet,
				Exec: func(ctx context.Context, args []string) error {
					var entry db.CreateEntryParams
					resp := &db.CreateEntryRow{}

					if len(args) != 1 {
						return fmt.Errorf("missing file name")
					}
					header, err := cromp.ParseFileHeader(args[0])
					if err != nil {
						return err
					}

					s := header.UUID.String()
					if s == "" || s == "00000000-0000-0000-0000-000000000000" {

						header.UUID = uuid.New()
					}

					entry.Title = header.Title
					entry.EntryID = header.UUID

					fmt.Printf("%#v\n", entry)

					data, err := cromp.ReadFileBody(args[0])
					if err != nil {
						return err
					}

					entry.Body = string(data)

					err = Post("/entries/add", entry, resp)
					if err != nil {
						return err
					}

					fmt.Printf("%#v\n", resp)

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
