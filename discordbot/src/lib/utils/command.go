package utils

import "github.com/bwmarrin/discordgo"

type Command struct {
	ApplicationCommand    *discordgo.ApplicationCommand
	DefaultHandler        func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
	ModalSubmittedHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
	ComponentHandler      func(session *discordgo.Session, interaction *discordgo.InteractionCreate, customId string)
}
