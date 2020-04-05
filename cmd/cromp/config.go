package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os/user"
	"path"

	"github.com/peterbourgon/ff/v2/ffcli"
)

func getPath() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	p := path.Join(usr.HomeDir, ".cromprc")

	return &p, nil
}

// Config represents the client configuration for cromp
type Config struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// WriteConfig dumbs the configuration to disk
func (c *Config) WriteConfig() error {

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	p, err := getPath()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*p, b, 0600)
}

// ReadConfig reads a configuration from disk
func (c *Config) ReadConfig() error {
	p, err := getPath()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(*p)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

// NewConfig creates a new config ffcli command
func NewConfig() *ffcli.Command {
	var configFlagSet = flag.NewFlagSet("cromp config", flag.ExitOnError)
	urlConfFS := configFlagSet.String("url", "", "URL of cromp server")
	tokenConfFS := configFlagSet.String("token", "", "Access token for cromp server")

	return &ffcli.Command{
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
	}
}
