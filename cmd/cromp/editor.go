package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
	cromp "suah.dev/cromp/internal"
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

			fmt.Printf("%#v\n", updatedEntry)

			return nil
		},
	}
}

// NewDoc opens a new doc in $EDITOR
func NewDoc() (*uuid.UUID, error) {
	var s string
	entry := &db.CreateEntryParams{}
	editor, args := editorCmd()

	entry.EntryID = uuid.New()

	f, err := ioutil.TempFile("", fmt.Sprintf("cromp-%s", entry.EntryID.String()))
	if err != nil {
		return nil, err
	}

	args = append(args, f.Name())

	defer os.Remove(f.Name())

	// TODO populate with a user defined template here
	f.Write([]byte(entry.Body))

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}
	s = string(b)

	header, err := cromp.ParseHeader(s)
	if err != nil {
		return nil, err
	}

	entry.Title = header.Title
	entry.Body = s

	result := make(map[string]string)

	err = Post("/entries/add", entry, result)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", result)

	return &entry.EntryID, nil
}

func editorCmd() (string, []string) {
	editor := "vi"
	if c := os.Getenv("EDITOR"); c != "" {
		editor = c
	}

	var args []string
	if pos := strings.Index(editor, " "); pos != -1 {
		args = strings.Fields(editor[pos+1:])
		editor = editor[:pos]
	}

	return editor, args
}

// OpenInEditor opens an existing cromp doc in $EDITOR
func OpenInEditor(id uuid.UUID) (*string, error) {
	var err error
	var params db.GetEntryParams
	var saveParams db.UpdateEntryParams

	s := ""
	resp := &db.Entry{}

	params.EntryID = id

	err = Post("/entries/get", params, resp)
	if err != nil {
		return nil, err
	}

	cmd, args := editorCmd()

	if f, err := ioutil.TempFile("", fmt.Sprintf("cromp-%s", resp.EntryID.String())); err == nil {
		args = append(args, f.Name())

		defer os.Remove(f.Name())
		f.Write([]byte(resp.Body))

		cmd := exec.Command(cmd, args...)
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

	header, err := cromp.ParseHeader(s)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", header)

	saveParams.EntryID = params.EntryID
	saveParams.Body = s
	saveParams.Title = header.Title

	result := make(map[string]string)

	err = Post("/entries/update", saveParams, result)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", result)

	return nil, err
}
