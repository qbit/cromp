package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
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

	done := make(chan bool)
	go func() {
		err = StartWatcher(f.Name(), *entry, saveEntry)
		if err != nil {
			log.Println(err)
		}
		<-done
	}()

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

	close(done)

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

// StartWatcher executes the callbacke "f" when the file is modified
func StartWatcher(p string, params interface{}, f func(string, interface{}) error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				err := f(p, params)
				if err != nil {
					log.Println(err)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()

	err = watcher.Add(p)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func saveEntry(fn string, p interface{}) error {
	var err error
	var createParams db.CreateEntryParams
	var updateParams db.UpdateEntryParams

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	s := string(b)

	header, err := cromp.ParseHeader(s)
	if err != nil {
		return err
	}

	result := make(map[string]string)

	switch p.(type) {
	case db.GetEntryParams:
		a, ok := p.(db.GetEntryParams)
		if !ok {
			return fmt.Errorf("invaled entry type")
		}
		updateParams.EntryID = a.EntryID
		updateParams.Body = s
		updateParams.Title = header.Title

		err = Post("/entries/update", updateParams, result)
		if err != nil {
			return err
		}
	case db.CreateEntryParams:
		fmt.Println("Create")
		a, ok := p.(db.CreateEntryParams)
		if !ok {
			return fmt.Errorf("invaled entry type")
		}
		createParams.EntryID = a.EntryID
		createParams.Body = s
		createParams.Title = header.Title

		err = Post("/entries/add", createParams, result)
		if err != nil {
			return err
		}
	}

	return nil
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

		done := make(chan bool)
		go func() {
			err = StartWatcher(f.Name(), params, saveEntry)
			if err != nil {
				log.Println(err)
			}
			<-done
		}()

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
		close(done)
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
