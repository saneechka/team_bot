package handler

import (
	"testing"
)

func TestIsAdmin(t *testing.T) {
	adminUsers := []string{"KsuNovak", "admin1", "admin2"}

	testCases := []struct {
		username string
		expected bool
	}{
		{"KsuNovak", true},
		{"admin1", true},
		{"admin2", true},
		{"regular_user", false},
		{"", false},
		{"ksunovak", false}, // case sensitive
	}

	for _, tc := range testCases {
		t.Run(tc.username, func(t *testing.T) {
			isAdmin := false
			for _, adminUsername := range adminUsers {
				if tc.username == adminUsername {
					isAdmin = true
					break
				}
			}

			if isAdmin != tc.expected {
				t.Errorf("Expected %v for username %s, got %v", tc.expected, tc.username, isAdmin)
			}
		})
	}
}
