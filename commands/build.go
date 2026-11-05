package commands

import (
	"github.com/keyvie/modmail/internal/dis"
	"github.com/keyvie/modmail/internal/state"
)

func Build() map[dis.ApplicationCommandKey]*dis.ApplicationCommand[state.Public] {
	commands := []*dis.ApplicationCommand[state.Public]{
		NewPing(),
	}

	all := make(map[dis.ApplicationCommandKey]*dis.ApplicationCommand[state.Public], len(commands))
	for _, command := range commands {
		all[command.Key()] = command
	}
	return all
}
