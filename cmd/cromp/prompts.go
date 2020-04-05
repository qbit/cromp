package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// SecurePrompt presents the user with a non-echoing prompt
func SecurePrompt(prompt string) (*string, error) {
	fmt.Print(prompt)
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	pass := string(b)
	return &pass, nil
}

// Prompt presents the user with an echoing prompt
func Prompt(prompt string) (*string, error) {
	var user string
	fmt.Print(prompt)
	fmt.Scanln(&user)
	return &user, nil
}
