package client

import "testing"

func TestAuthorizationHeaderValue(t *testing.T) {
	tests := []struct {
		name      string
		basicAuth string
		token     string
		want      string
	}{
		{
			name:      "token takes priority when both are set",
			basicAuth: "dXNlcjpwYXNz",
			token:     "abc123",
			want:      "Bearer abc123",
		},
		{
			name:      "basic auth used when no token",
			basicAuth: "dXNlcjpwYXNz",
			token:     "",
			want:      "Basic dXNlcjpwYXNz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authorizationHeaderValue(tt.basicAuth, tt.token)
			if got != tt.want {
				t.Errorf("authorizationHeaderValue(%q, %q) = %q, want %q", tt.basicAuth, tt.token, got, tt.want)
			}
		})
	}
}
