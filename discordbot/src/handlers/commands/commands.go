package commands

import (
	"github.com/bwmarrin/discordgo"

	"discordbot/src/handlers/commands/minecraft"
)

var Commands = map[*discordgo.ApplicationCommand]func(*discordgo.Session, *discordgo.InteractionCreate){
	minecraft.Command: minecraft.Handler,
}
