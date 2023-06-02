package minecraft

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"

	"discordbot/src/lib/colours"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "minecraft",
	Description: "Creates a message that contains all the details about the minecraft server.",
}

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Member != nil && interaction.Member.User.ID == os.Getenv("OWNER_ID") {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Bedrock Connection Details",
								Value: fmt.Sprintf("Server address: `%s`\nPort: `%s`\n​", os.Getenv("BEDROCK_ADDRESS"), os.Getenv("BEDROCK_PORT")),
							},
							{
								Name:  "Java Connection Details",
								Value: fmt.Sprintf("Server address: `%s:%s`\n​", os.Getenv("JAVA_ADDRESS"), os.Getenv("JAVA_PORT")),
							},
							{
								Name:  "Worlds",
								Value: "Hub (access via `/hub`)\nNew world (access via `/newworld`)\nOld world (access via `/oldworld`)\n​",
							},
							{
								Name:  "Server Status",
								Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Offline`\n`Users: None`\n",
							},
						},
						Color: colours.ColourRed,
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Start server",
								Style:    discordgo.SuccessButton,
								Disabled: false,
								CustomID: "minecraft:start",
							},
							discordgo.Button{
								Label:    "Stop server",
								Style:    discordgo.DangerButton,
								Disabled: false,
								CustomID: "minecraft:stop",
							},
						},
					},
				},
			},
		})
	} else {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "You don't have the required permissions to run this command.",
						Color:       colours.ColourBlue,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
