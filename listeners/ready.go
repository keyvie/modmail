package listeners

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/keyvie/modmail/internal/dis"
	"github.com/keyvie/modmail/internal/state"
)

func NewReady() *dis.Listener[state.Public] {
	return dis.NewListener(func(container *dis.Container[state.Public], event *events.Ready) error {
		slog.Info(fmt.Sprintf("Logged in as %s#%s (%s)", event.User.Username, event.User.Discriminator, event.User.ID))

		commands := make([]discord.ApplicationCommandCreate, 0, len(container.ApplicationCommands))
		for _, command := range container.ApplicationCommands {
			commands = append(commands, command.Data)
		}
		_, err := container.Client.Rest.SetGlobalCommands(event.User.ID, commands)
		if err != nil {
			return fmt.Errorf("failed to register global commands: %w", err)
		}

		return nil
	})
}
