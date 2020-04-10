package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
)

// NewEditor makes the "cromp edit" command
func NewEditor() *ffcli.Command {
	var editorFlagSet = flag.NewFlagSet("cromp edit", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "edit",
		ShortUsage: "cromp edit id",
		FlagSet:    editorFlagSet,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("please specify the document id")
			}
			id, err := uuid.Parse(args[0])
			if err != nil {
				return err
			}
			updatedEntry, err := OpenInEditor(id)
			if err != nil {
				return err
			}

			fmt.Println(*updatedEntry)

			return nil
		},
	}
}

// OpenInEditor opens a cromp doc in $EDITOR
func OpenInEditor(id uuid.UUID) (*string, error) {
	var err error
	var params db.GetEntryParams
	resp := &db.Entry{}
	editor := "vi"
	s := ""

	params.EntryID = id

	if c := os.Getenv("EDITOR"); c != "" {
		editor = c
	}
	err = Post("/entries/get", params, resp)
	if err != nil {
		return nil, err
	}

	if f, err := ioutil.TempFile("", fmt.Sprintf("cromp-%s", resp.EntryID.String())); err == nil {
		defer os.Remove(f.Name())
		f.Write([]byte(resp.Body))
		cmd := exec.Command(editor, f.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadFile(f.Name())
		if err != nil {
			return nil, err
		}
		s = string(b)
	}

	return &s, err
}
