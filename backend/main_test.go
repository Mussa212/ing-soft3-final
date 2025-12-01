package main

import "testing"

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
		name: "masks password",
		dsn:  "postgres://user:secret@localhost:5432/db?sslmode=disable",
		want: "postgres://user:%2A%2A%2A%2A@localhost:5432/db?sslmode=disable",
	},
		{
			name: "returns original when no credentials",
			dsn:  "postgres://localhost:5432/db",
			want: "postgres://localhost:5432/db",
		},
		{
			name: "returns original when parse fails",
			dsn:  "not a valid dsn %",
			want: "not a valid dsn %",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactDSN(tt.dsn); got != tt.want {
				t.Fatalf("redactDSN() = %s, want %s", got, tt.want)
			}
		})
	}
}
