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

// ReadPassword reads a password and confirmation from the terminal.
func ReadPassword() (string, error) {
	for i := 0; i < maxPromptRetries; i++ {
		pass, err := readPass("Password: ")
		if err != nil {
			return "", err
		}

		confirmPass, err := readPass("Confirm password: ")
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

func readPass(prompt string) (string, error) {
	fmt.Print(prompt)
	b, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Print("\n")
	return string(b), err
}
