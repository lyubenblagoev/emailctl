package emailctl

import (
	"github.com/lyubenblagoev/goprsc"
)

// AuthResponse is a wrapper for goprsc.AuthResponse
type AuthResponse struct {
	*goprsc.AuthResponse
}

// AuthService handles communication with the authentication API of the Postfix REST Server.
type AuthService service

// Login authenticates the user credential and returnes the tokens provided by the Postfix REST Server.
func (s *AuthService) Login(login, password string) (*AuthResponse, error) {
	response, err := s.client.Auth.Login(login, password)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{AuthResponse: response}, nil
}

// Logout logs out the user, if the provided login and refreshToken are valid.
func (s *AuthService) Logout(login, refreshToken string) error {
	err := s.client.Auth.Logout(login, refreshToken)
	return err
}
