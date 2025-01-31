package cli

import (
	"fmt"
	"os"
)

// VerbHandler is a function that handles a command verb
type VerbHandler func(args []string) error

// HandlerMap is a type alias for a map of command handlers
type HandlerMap map[string]VerbHandler

// handlers is the global map of verb handlers
var handlers = HandlerMap{}

// HandleVerb processes the command line arguments and executes the appropriate handler
func HandleVerb() error {
	if _, exists := handlers["help"]; !exists {
		return fmt.Errorf("help command not found")
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("no command provided. Usage: %s <command> [args...]", os.Args[0])
	}

	verb := os.Args[1]
	handler, exists := handlers[verb]
	if !exists {
		return fmt.Errorf("unknown command '%s'. Use '%s help' to see available commands", 
			verb, os.Args[0])
	}

	return handler(os.Args[2:])
}

// AddVerbHandler adds a verb handler to the map
func AddVerbHandler(verb string, handler VerbHandler) error {
	if verb == "" {
		return fmt.Errorf("verb cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	if _, exists := handlers[verb]; exists {
		return fmt.Errorf("handler for verb '%s' already registered", verb)
	}

	handlers[verb] = handler
	return nil
}

// ListVerbs returns a slice of all registered verb names
func ListVerbs() []string {
	verbs := make([]string, 0, len(handlers))
	for verb := range handlers {
		verbs = append(verbs, verb)
	}
	return verbs
}
