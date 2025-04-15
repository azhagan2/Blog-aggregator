package state

import (
	"github.com/azhagan2/blog_aggregator/internal/config"
	"github.com/azhagan2/blog_aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

// constructor

func New(cfg *config.Config, dbQueries *database.Queries) *State {
	return &State{
		Cfg: cfg,
		Db:  dbQueries,
	}
}

/* what is happening is, when I create a State struct with cfg *Config (the prefix is just convention to follow for
Exporting) Then It creates a pointer to Config struct in the config.go file. Then using a New() function with parameter of
cfg, a real instance (or value) of Config, being assigned to the Cfg of State struct variable.

So, from this I can observe that, the State struct acts as a container or context that acts as the same which is
being referred to, May be this is for encapsulation ?

What you're seeing is less about hiding data (traditional encapsulation) and more about organizing related dependencies.
It's a form of composition where you build a new structure that contains other structures.*/

// State is a container for application-wide dependencies and shared resources.
// It follows the dependency injection pattern, allowing command handlers to access
// common resources (initially just configuration, but potentially databases, loggers,
// or other services in the future) through a single parameter rather than multiple ones.
//
// By centralizing dependencies in this struct:
// 1. Command handler signatures remain consistent even as dependencies change
// 2. Testing becomes easier through mockable state objects
// 3. Dependencies are explicitly managed rather than using globals
// 4. New shared resources can be added without modifying function signatures
//
// The State struct serves as the context in which commands execute, giving them
// access to everything they need to perform their functions.
