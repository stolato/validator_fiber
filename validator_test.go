package validator_fiber_test

import (
	"testing"

	"github.com/stolato/validator_fiber"
)

// User struct para testar a validação
type User struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Age      uint   `json:"age" validate:"required,gte=18"`
	Password string `json:"password" validate:"required,min=6"`
}

func TestValidator(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		expectedError  bool
		expectedErrors []struct {
			Field string
			Tag   string
		}
	}{
		{
			name: "Valid User - No Errors",
			input: User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Age:      25,
				Password: "securepassword",
			},
			expectedError: false,
		},
		{
			name: "Invalid User - Missing Name",
			input: User{
				Email:    "jane.doe@example.com",
				Age:      30,
				Password: "anothersecurepassword",
			},
			expectedError: true,
			expectedErrors: []struct {
				Field string
				Tag   string
			}{
				{Field: "name", Tag: "required"},
			},
		},
		{
			name: "Invalid User - Invalid Email and Short Password",
			input: User{
				Name:     "Alice",
				Email:    "invalid-email",
				Age:      22,
				Password: "short",
			},
			expectedError: true,
			expectedErrors: []struct {
				Field string
				Tag   string
			}{
				{Field: "email", Tag: "email"},
				{Field: "password", Tag: "min"},
			},
		},
		{
			name: "Invalid User - Age Too Young",
			input: User{
				Name:     "Bob",
				Email:    "bob@example.com",
				Age:      17,
				Password: "bobs_password",
			},
			expectedError: true,
			expectedErrors: []struct {
				Field string
				Tag   string
			}{
				{Field: "age", Tag: "gte"},
			},
		},
		{
			name: "Invalid User - Name Too Short",
			input: User{
				Name:     "A",
				Email:    "charlie@example.com",
				Age:      28,
				Password: "charlies_password",
			},
			expectedError: true,
			expectedErrors: []struct {
				Field string
				Tag   string
			}{
				{Field: "name", Tag: "min"},
			},
		},
		{
			name: "Invalid User - Multiple Errors",
			input: User{
				Name:     "",             // Missing
				Email:    "not-an-email", // Invalid
				Age:      10,             // Too young
				Password: "abc",          // Too short
			},
			expectedError: true,
			expectedErrors: []struct {
				Field string
				Tag   string
			}{
				{Field: "name", Tag: "required"},
				{Field: "email", Tag: "email"},
				{Field: "age", Tag: "gte"},
				{Field: "password", Tag: "min"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator_fiber.Validator(tt.input)

			if tt.expectedError {
				if len(errs) == 0 {
					t.Errorf("Expected errors but got none for test: %s", tt.name)
				}
				if len(errs) != len(tt.expectedErrors) {
					t.Errorf("Expected %d errors, got %d for test: %s", len(tt.expectedErrors), len(errs), tt.name)
				}

				// Check if the expected errors are present and correct
				for _, expectedErr := range tt.expectedErrors {
					found := false
					for _, actualErr := range errs {
						if actualErr.FailedField == expectedErr.Field && actualErr.Tag == expectedErr.Tag {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error for field '%s' with tag '%s' not found for test: %s", expectedErr.Field, expectedErr.Tag, tt.name)
					}
				}

			} else {
				if len(errs) > 0 {
					t.Errorf("Expected no errors but got %d errors for test: %s", len(errs), tt.name)
					for _, err := range errs {
						t.Logf("  Error: Field: %s, Tag: %s, Value: %v", err.FailedField, err.Tag, err.Value)
					}
				}
			}
		})
	}
}
