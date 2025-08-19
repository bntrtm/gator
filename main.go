package main

import (
	_ "github.com/lib/pq"
	"fmt"
	"os"
	"time"
	"database/sql"

	"github.com/bntrtm/gator/internal/config"
	"github.com/bntrtm/gator/internal/database"
)

func main() {
	cliState := state{}
	cfg := config.Read()
	cliState.config = &cfg
	db, err := sql.Open("postgres", cliState.config.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	cliState.db = dbQueries

	cmdRegistry := commands{
		handlers: make(map[string]func(*state, command) error),
	
	}
	cmdRegistry.register("reset", handlerReset)
	cmdRegistry.register("register", handlerRegister)
	cmdRegistry.register("login", handlerLogin)
	cmdRegistry.register("users", handlerUsers)
	cmdRegistry.register("agg", handlerAgg)
	cmdRegistry.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmdRegistry.register("feeds", handlerFeeds)
	cmdRegistry.register("follow", middlewareLoggedIn(handlerFollow))
	cmdRegistry.register("following", middlewareLoggedIn(handlerFollowing))
	
	client := NewClient(time.Second * 10)
	cliState.client = &client

	args := os.Args
	if len(args) < 2 {
		fmt.Println("ERROR: Not enough arguments")
		os.Exit(1)
	}
	args = args[1:]
	command := command{name: args[0], args: args[1:]}
	err = cmdRegistry.run(&cliState, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}