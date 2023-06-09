package minecraft

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordbot/src/lib/colours"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "minecraft",
	Description: "Creates a message that contains all the details about the minecraft server.",
}

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Member != nil && interaction.Member.User.ID == os.Getenv("OWNER_ID") {
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Bedrock Connection Details",
				Value: fmt.Sprintf("Server address: `%s`\nPort: `%s`\n​", os.Getenv("BEDROCK_ADDRESS"), os.Getenv("BEDROCK_PORT")),
			},
			{
				Name:  "Java Connection Details",
				Value: fmt.Sprintf("Server address: `%s:%s`\n​", os.Getenv("JAVA_ADDRESS"), os.Getenv("JAVA_PORT")),
			},
		}

		additional := os.Getenv("ADDITIONAL_MESSAGES_FOR_EMBED")
		for _, field := range strings.Split(additional, ";;") {
			fieldParts := strings.Split(field, "::")
			if len(fieldParts) == 2 {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  fieldParts[0],
					Value: strings.Replace(fieldParts[1], `\n`, "\n", -1) + "\n​",
				})
			}
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Server Version",
			Value: fmt.Sprintf("The server is currently running version `%s`\n​", os.Getenv("MC_VERSION")),
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Server Status",
			Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Offline`\n`Users: None`\n",
		})

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Fields: fields,
						Color:  colours.ColourRed,
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
