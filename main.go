package main

import (
	"fmt"
	"os"

	"github.com/bntrtm/gator/internal/config"
)

func main() {
	cliState := state{}
	cfg := config.Read()
	cliState.config = &cfg
	cmdRegistry := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmdRegistry.register("login", handlerLogin)
	
	args := os.Args
	if len(args) < 2 {
		fmt.Println("ERROR: Not enough arguments")
		os.Exit(1)
	}
	args = args[1:]
	command := command{name: args[0], args: args[1:]}
	err := cmdRegistry.run(&cliState, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}