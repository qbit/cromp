package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
	cromp "suah.dev/cromp/internal"
)

// NewLoad returns the load ffcli command
func NewLoad() *ffcli.Command {
	return &ffcli.Command{
		Name:       "load",
		ShortUsage: "cromp load <file>",
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
	}
}
