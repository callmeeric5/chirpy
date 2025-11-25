package auth_test

import (
	"go_server/internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        string
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("MakeJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
