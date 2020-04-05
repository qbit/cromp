package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
)

// NewSimilar creates a new config ffcli command
func NewSimilar() *ffcli.Command {
	return &ffcli.Command{
		Name:       "similar",
		ShortUsage: "cromp similar [text]",
		Exec: func(ctx context.Context, args []string) error {
			var params db.SimilarEntriesParams
			resp := &[]db.SimilarEntriesRow{}

			if len(args) != 1 {
				return fmt.Errorf("missing file name")
			}

			pos, err := GetPOS(args[0])
			if err != nil {
				return err
			}

			params.Similarity = strings.Join(pos["NN"], "|")

			err = Post("/entries/similar", params, resp)
			if err != nil {
				return err
			}

			for _, e := range *resp {
				fmt.Printf("%s\t%s\n",
					e.Title,
					e.Headline)
			}

			return nil
		},
	}
}
