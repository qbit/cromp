package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
	"gopkg.in/jdkato/prose.v2"
)

// GetPOS returns the parts of speech from a text file
func GetPOS(p string) (map[string][]string, error) {
	// We only care about word like things
	var wre = regexp.MustCompile(`^\w+$`)
	var pos = map[string][]string{}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	doc, err := prose.NewDocument(string(data))
	if err != nil {
		return nil, err
	}

	for _, tok := range doc.Tokens() {
		if wre.MatchString(tok.Text) {
			pos[tok.Tag] = append(pos[tok.Tag], tok.Text)
		}
	}

	return pos, nil
}

// NewPOS returns a new ffcli.Command
func NewPOS() *ffcli.Command {
	return &ffcli.Command{
		Name:       "pos",
		ShortUsage: "cromp pos <file>",
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing file name")
			}

			pos, err := GetPOS(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Nouns: %s\n", strings.Join(pos["NN"], ", "))
			fmt.Printf("Verbs: %s\n", strings.Join(pos["VB"], ", "))
			return nil
		},
	}
}
