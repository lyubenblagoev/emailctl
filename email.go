package emailctl

import (
	"fmt"
	"net/mail"
)

// ValidateEmailFromParts validates the email address using the given username
// and domain and returns an error if the email address is invalid.
func ValidateEmailFromParts(name, domain string) error {
	email := fmt.Sprintf("%s@%s", name, domain)
	return ValidateEmail(email)
}

// ValidateEmail validates the given email address and returns an error if
// the email address is invalid.
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: '%s'", email)
	}
	return nil
}
