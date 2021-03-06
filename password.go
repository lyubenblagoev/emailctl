package emailctl

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	maxPromptRetries = 3
)

// ReadAndConfirmPassword reads a password and confirmation from the terminal.
// Retries three times if the passwords do not match.
func ReadAndConfirmPassword() (string, error) {
	for i := 0; i < maxPromptRetries; i++ {
		pass, err := ReadPassword("Password: ")
		if err != nil {
			return "", err
		}

		confirmPass, err := ReadPassword("Confirm password: ")
		if err != nil {
			return "", err
		}

		if strings.Compare(pass, confirmPass) != 0 {
			if i < maxPromptRetries {
				fmt.Println("Passwords don't match!")
			}
			continue
		}

		return pass, nil
	}

	return "", errors.New("Passwords don't match")
}

// ReadPassword reads a password from the terminal
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	b, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Print("\n")
	return string(b), err
}
