package stop

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"discordbot/src/lib/colours"
	"discordbot/src/lib/timeout"
)

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if !timeout.GetTimeout("minecraft") {
		durationInSeconds, err := strconv.Atoi(os.Getenv("START_STOP_TIMEOUT_IN_SECONDS"))
		if err != nil {
			durationInSeconds = 30
		}
		go timeout.StartTimeout("minecraft", time.Second*time.Duration(durationInSeconds))

		session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      interaction.Message.ID,
			Channel: interaction.Message.ChannelID,
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
		})

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Stopping minecraft server...",
						Color:       colours.ColourBlue,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
			Description: fmt.Sprintf("%s has stopped the server", interaction.Member.Mention()),
			Color:       colours.ColourRed,
		})

		cmd := exec.Command("pkill", "java")
		cmd.Start()
	} else {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Please wait a short period of time before using this action again",
						Color:       colours.ColourBlue,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
