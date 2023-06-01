package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"discordbot/src/handlers/commands"
	"discordbot/src/handlers/components"
)

func CommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionMessageComponent:
		for component, handler := range components.Components {
			if component == interaction.MessageComponentData().CustomID {
				log.Printf("[%s]: Message component %s\n", interaction.Member.User.String(), component)
				handler(session, interaction)
			}
		}
	default:
		for command, handler := range commands.Commands {
			if command.Name == interaction.ApplicationCommandData().Name {
				log.Printf("[%s]: /%s\n", interaction.Member.User.String(), command.Name)
				handler(session, interaction)
			}
		}
	}
}

func RegisterCommands(session *discordgo.Session) {
	for command := range commands.Commands {
		session.ApplicationCommandCreate(os.Getenv("APP_ID"), os.Getenv("GUILD_ID"), command)
	}
}

func UnregisterCommands(session *discordgo.Session) {
	for command := range commands.Commands {
		session.ApplicationCommandDelete(os.Getenv("APP_ID"), os.Getenv("GUILD_ID"), command.ID)
	}
}

func init() {
	godotenv.Load()
}

func main() {
	session, sessionError := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if sessionError != nil {
		panic(sessionError)
	}

	fmt.Println("Starting bot...")

	session.AddHandler(CommandHandler)
	RegisterCommands(session)

	session.Open()

	fmt.Println("Bot started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Closing bot gracefully...")

	UnregisterCommands(session)

	session.Close()

	fmt.Println("Bot closed gracefully")
}
