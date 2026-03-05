package validator

import (
	"testing"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestNew(t *testing.T) {
	v, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if v == nil {
		t.Error("New() returned nil validator")
	}
}

func TestValidate(t *testing.T) {
	v, _ := New()

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "valid input",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			input: TestStruct{
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:  "John",
				Email: "invalid-email",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "age out of range",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   200,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVar(t *testing.T) {
	v, _ := New()

	tests := []struct {
		name    string
		value   interface{}
		tag     string
		wantErr bool
	}{
		{"required string", "test", "required", false},
		{"empty string required", "", "required", true},
		{"email valid", "test@example.com", "email", false},
		{"email invalid", "invalid", "email", true},
		{"min length", "ab", "min=3", true},
		{"max length", "abc", "max=3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.value, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Var() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
