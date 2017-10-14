package emailctl

import (
	"fmt"
	"net/mail"
)

// ValidateEmailAddress validates the email address using the given account name and domain
// and returns an error if the name is invalid.
func ValidateEmailAddress(name, domain string) error {
	email := fmt.Sprintf("%s@%s", name, domain)
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: '%s'; %v", email, err)
	}
	return nil
}
