package cli

import "os"

// verbHandler is a function that handles a command verb
type verbHandler func(verb string, args []string)

// verbHandlers is a map of verb handlers
var verbHandlers = map[string]verbHandler{}

func HandleVerb(verbHandler verbHandler) {
	verbHandler(os.Args[1], os.Args[2:])
}

func AddVerbHandler(verb string, handler verbHandler) {
	verbHandlers[verb] = handler
}
