package mocks

import "snippetbox.brainwhat/internal/models"

// These are here for testing purposes

type MockUserModel struct{}

func (m *MockUserModel) Create(email, name, password string) error {
	switch email {
	case "dupe@mail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *MockUserModel) Authenticate(email, password string) (int, error) {
	if email == "example@mail.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *MockUserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
