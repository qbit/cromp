package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get is a generic HTTP GET handler
func Get(url string, resBody interface{}) error {
	var req *http.Request
	client := http.DefaultClient
	buf := new(bytes.Buffer)

	cfg := &Config{}
	err := cfg.ReadConfig()
	if err != nil {
		return err
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("%s%s", cfg.URL, url), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", cfg.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		if err = json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Println(string(data))

			return err
		}
	} else {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
	}

	return nil
}

// Post is a generic HTTP POST handler
func Post(url string, reqBody, resBody interface{}) (err error) {
	var req *http.Request
	client := http.DefaultClient
	buf := new(bytes.Buffer)

	cfg := &Config{}
	err = cfg.ReadConfig()
	if err != nil {
		return err
	}

	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return err
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s%s", cfg.URL, url), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", cfg.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 199 && res.StatusCode < 300 {
		if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return err
		}
	} else {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
	}

	return nil
}
