package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pl1000100/gator/internal/config"
	"github.com/pl1000100/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	s := state{
		cfg: &cfg,
	}

	db, err := sql.Open("postgres", s.cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	commands := commands{
		NameFunc: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)

	if len(os.Args) < 2 {
		log.Fatalf("error not enough arguments")
	}
	command := command{
		Name: os.Args[1],
	}
	if len(os.Args) > 2 {
		command.Arguments = os.Args[2:]
	}

	commands.run(&s, command)

}
