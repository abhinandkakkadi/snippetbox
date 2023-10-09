package mocks

import "github.com/abhinandkakkadi/snippetbox/internal/models"


type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pas$$word" {
		return 1, nil
	}
	
	return 0, models.ErrInvalidCredentials
}