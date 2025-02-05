package mocks

import (
	"time"

	"snippetbox.brainwhat/internal/models"
)

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
func (m *MockUserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		u := &models.User{
			Name:    "Ivan Ivanov",
			Email:   "example@mail.com",
			Created: time.Now(),
		}
		return u, nil
	} else {
		return nil, models.ErrNoRecord
	}
}
