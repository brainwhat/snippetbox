package mocks

import (
	"time"

	"snippetbox.brainwhat/internal/models"
)

// These are here for testing purposes

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type MockSnippetModel struct{}

func (m *MockSnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

func (m *MockSnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *MockSnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
