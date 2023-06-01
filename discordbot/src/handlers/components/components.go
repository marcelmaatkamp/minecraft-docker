package components

import (
	"github.com/bwmarrin/discordgo"

	minecraftstart "discordbot/src/handlers/components/minecraft/start"
	minecraftstop "discordbot/src/handlers/components/minecraft/stop"
)

var Components = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"minecraft:start": minecraftstart.Handler,
	"minecraft:stop":  minecraftstop.Handler,
}
