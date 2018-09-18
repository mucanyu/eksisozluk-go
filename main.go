package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	// cmds "github.com/mucanyu/eksisozluk-go/commands"
)

func main() {
	c := cli.NewCLI("eksisozluk-go", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"gundem": func() (cli.Command, error) {
			return &Gundem{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
