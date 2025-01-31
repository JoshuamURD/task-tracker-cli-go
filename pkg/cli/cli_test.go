package cli

import (
	"fmt"
	"os"
	"testing"
)

func TestHandleVerb(t *testing.T) {
	test := []struct {
		name string
		args []string
		want error
	}{
		{name: "no command", args: []string{}, want: fmt.Errorf("no command provided. Usage: %s <command> [args...]", os.Args[0])},
	}
	
}

func TestAddVerbHandler(t *testing.T) {
	// Reset handlers map before all tests
	handlers = map[string]VerbHandler{}

	tests := []struct {
		name    string
		verb    string
		handler VerbHandler
		wantErr string // Expected error message
	}{
		{
			name:    "nil handler",
			verb:    "help",
			handler: nil,
			wantErr: "handler cannot be nil",
		},
		{
			name:    "empty verb",
			verb:    "",
			handler: nil,
			wantErr: "verb cannot be empty",
		},
		{
			name:    "valid handler",
			verb:    "help",
			handler: func(verb string, args []string) error { return nil },
			wantErr: "",
		},
		{
			name:    "duplicate verb",
			verb:    "help",
			handler: func(verb string, args []string) error { return nil },
			wantErr: "handler for verb 'help' already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddVerbHandler(tt.verb, tt.handler)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("AddVerbHandler() error = nil, wantErr %q", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr {
					t.Errorf("AddVerbHandler() error = %q, wantErr %q", err.Error(), tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("AddVerbHandler() unexpected error = %v", err)
			}
		})
	}
}
