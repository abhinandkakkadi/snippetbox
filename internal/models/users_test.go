package models

import (
	"testing"

	"github.com/abhinandkakkadi/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {

	// Set up a suite of table-driven tests and expected results
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		}, {
			name:   "Zero ID",
			userID: 0,
			want:   false,
		}, {
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the testDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test

			db := newTestDB(t)

			// Create a new instance fof the UserModel
			m := UserModel{DB: db}

			// call the UserModel.Exists() method and check that the return value
			// and error match the expected values of the sub-test
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
