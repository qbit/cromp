package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path"
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
