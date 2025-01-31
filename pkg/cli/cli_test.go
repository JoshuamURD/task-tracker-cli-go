package cli

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestHandleVerb(t *testing.T) {
	handlers = HandlerMap{
		"help": func(verb string, args []string) error { return nil },
	}
	test := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "no command",
			args:    []string{},
			wantErr: fmt.Sprintf("no command provided. Usage: %s <command> [args...]", os.Args[0]),
		},
		{
			name:    "help",
			args:    []string{"help"},
			wantErr: "",
		},
		{
			name:    "unknown command",
			args:    []string{"unknown"},
			wantErr: fmt.Sprintf("unknown command '%s'. Use '%s help' to see available commands", "unknown", os.Args[0]),
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append([]string{os.Args[0]}, tt.args...)
			got := HandleVerb()
			if tt.wantErr == "" {
				if got != nil {
					t.Errorf("Testing HandleVerb() by %s: Got unexpected error: %v", tt.name, got)
				}
			} else if got == nil {
				t.Errorf("Testing HandleVerb() by %s: Expected error %q but got nil", tt.name, tt.wantErr)
			} else if got.Error() != tt.wantErr {
				t.Errorf("Testing HandleVerb() by %s: Got %q, want %q", tt.name, got.Error(), tt.wantErr)
			}
		})
	}
}

func TestAddVerbHandler(t *testing.T) {
	// Reset handlers map before all tests
	handlers = HandlerMap{}

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

func TestListVerbs(t *testing.T) {
	handlers = HandlerMap{
		"help": func(verb string, args []string) error { return nil },
	}
	got := ListVerbs()
	want := []string{"help"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListVerbs() = %v, want %v", got, want)
	}
}
