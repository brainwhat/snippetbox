package models

import (
	"testing"

	"snippetbox.brainwhat/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userId int
		want   bool
	}{
		{
			name:   "ValidID",
			userId: 1,
			want:   true,
		},
		{
			name:   "FalseID",
			userId: 2,
			want:   false,
		},
		{
			name:   "ZeroID",
			userId: 0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(tt.userId)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}

}
