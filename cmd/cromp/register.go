package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
)

// NewRegister creates a new config ffcli command
func NewRegister() *ffcli.Command {
	var regFlagSet = flag.NewFlagSet("cromp reg", flag.ExitOnError)
	urlConfFS := regFlagSet.String("url", "", "URL of cromp server")
	return &ffcli.Command{
		Name:       "reg",
		ShortUsage: "reg -url [url]",
		FlagSet:    regFlagSet,
		Exec: func(ctx context.Context, args []string) error {
			var regParams db.CreateUserParams
			resp := &db.CreateUserRow{}

			cfg := &Config{}
			err := cfg.ReadConfig()
			if err != nil {
				return err
			}

			if cfg.URL == "" && *urlConfFS == "" {
				return fmt.Errorf("please specify -url")
			}

			fn, err := Prompt("First Name: ")
			if err != nil {
				return err
			}

			ln, err := Prompt("Last Name: ")
			if err != nil {
				return err
			}

			email, err := Prompt("Email: ")
			if err != nil {
				return err
			}

			user, err := Prompt("Username: ")
			if err != nil {
				return err
			}
			pass, err := SecurePrompt("Password: ")
			if err != nil {
				return err
			}

			fmt.Println()

			regParams.Username = *user
			regParams.Password = *pass
			regParams.FirstName = *fn
			regParams.LastName = *ln
			regParams.Email = *email

			err = Post("/user/new", regParams, resp)
			if err != nil {
				return err
			}

			fmt.Printf("New token expires %s\n", resp.TokenExpires)

			cfg.Token = resp.Token
			err = cfg.WriteConfig()
			if err != nil {
				return err
			}

			return nil
		},
	}
}
