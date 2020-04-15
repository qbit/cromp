package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v2/ffcli"
	"suah.dev/cromp/db"
)

// NewLogin creates a new config ffcli command
func NewLogin() *ffcli.Command {
	var loginFlagSet = flag.NewFlagSet("cromp login", flag.ExitOnError)
	urlConfFS := loginFlagSet.String("url", "", "URL of cromp server")
	return &ffcli.Command{
		Name:       "login",
		ShortUsage: "login -url [url]",
		FlagSet:    loginFlagSet,
		Exec: func(ctx context.Context, args []string) error {
			var loginParams db.AuthUserParams
			resp := &db.AuthUserRow{}

			cfg := &Config{}
			err := cfg.ReadConfig()
			if err != nil {
				return err
			}

			if cfg.URL == "" && *urlConfFS == "" {
				return fmt.Errorf("please specify -url")
			}

			cfg.URL = *urlConfFS
			err = cfg.WriteConfig()
			if err != nil {
				return err
			}

			user, err := Prompt("Login: ")
			if err != nil {
				return err
			}
			pass, err := SecurePrompt("Password: ")
			if err != nil {
				return err
			}

			fmt.Println()

			loginParams.Username = *user
			loginParams.Crypt = *pass

			err = Post("/user/auth", loginParams, resp)
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
