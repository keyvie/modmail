package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"

	"github.com/joho/godotenv"

	"github.com/keyvie/modmail/commands"
	"github.com/keyvie/modmail/internal/database"
	"github.com/keyvie/modmail/internal/dis"
	"github.com/keyvie/modmail/internal/state"
	"github.com/keyvie/modmail/listeners"
)

func main() {
	// Load env values from .env if it exists
	_ = godotenv.Load()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// Listen for signals to shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Checks for database connection url
	mongoUri, exists := os.LookupEnv("MONGO_URI")
	if !exists || mongoUri == "" {
		panic("MONGO_URI is missing")
	}

	data, err := database.NewDatabase(mongoUri)
	if err != nil {
		panic(err)
	}

	// Open the discord client
	discordClient, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
				gateway.IntentMessageContent,
			),
		),
	)
	if err != nil {
		panic(err)
	}

	container := dis.NewContainer(
		discordClient,
		dis.WithState(*state.New(data)),
		dis.WithListeners(listeners.Build()),
		dis.WithApplicationCommands(commands.Build()),
	)
	dis.Initialize(container)

	// Open connection to Discord
	if err = discordClient.OpenGateway(ctx); err != nil {
		panic(err)
	}

	// What gets ran after shutdown is triggered
	<-ctx.Done()
	stop()

	// Make a context for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	discordClient.Close(shutdownCtx)

	if err := data.Conn.Disconnect(shutdownCtx); err != nil {
		slog.Error("Mongo", "Error closing", err)
	}
}
