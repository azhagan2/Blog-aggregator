package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/azhagan2/blog_aggregator/internal/command"
	"github.com/azhagan2/blog_aggregator/internal/config"
	"github.com/azhagan2/blog_aggregator/internal/database"
	"github.com/azhagan2/blog_aggregator/internal/state"
)

// This state struct is to desgin pattern, which enables us to add more sb connections later on.

/* This clicommand struct is nothing but, stores the "type of command" as name, followed by giving commands as "arguments"
for example, giving commands like, login, register, follow, these are commands, followed by giving additional arguments,
are stored in the argument slice of this struct */

func main() {

	// First using the Read func in the config.go, we read the config file, and store it in a GO struct.
	// The cfg basically contains the JSON converted Go struct (Config). {postgres://example, azhagan2}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading the json file", err)
		os.Exit(1)
	}

	dbURL := "postgres://postgres:postgres@localhost:5432/gator"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("error in opening a new connection,:", err)
		return
	}

	dbQueries := database.New(db)

	s := state.New(cfg, dbQueries)

	/* Then we assign a new user, where there is already a place in the skeleton of the Config struct, it will assigns
	and then inside this SetUser func, write func is called and it converts and store(write) it in the config file. */

	// cfg.SetUser("Juicy_jelly")

	/* Again for the testing sake, we are reading, as JSON -> GO STRUCT, then passing that struct here to display
	to prove we storing actually*/

	cmds := command.NewCommands()

	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandlerReset)
	cmds.Register("users", command.HandlerGetUsers)
	cmds.Register("agg", command.HandlerAgg)
	cmds.Register("addfeed", command.MiddlewareLoggedIn(command.HandlerAddfeed))
	cmds.Register("feeds", command.HandlerFeeds)
	cmds.Register("follow", command.MiddlewareLoggedIn(command.HandlerFollow))
	cmds.Register("following", command.MiddlewareLoggedIn(command.HandlerFollowing))
	cmds.Register("unfollow", command.MiddlewareLoggedIn(command.HandlerUnfollow))
	cmds.Register("browse", command.MiddlewareLoggedIn(command.HandlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}

	cmnd := command.Clicommand{
		Name:     os.Args[1],
		Argument: os.Args[2:],
	}

	/*This Run function call acts as a bridge between CLI interface (state s) and actual functionality (handler logic func)
	In fp perspective, the Run func, passing the First class func, giving resources they need to do the job*/

	if err := cmds.Run(s, cmnd); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// fmt.Println(cfg.DbURL)
	// fmt.Println(cfg.CurrentUserName)
}
