package stop

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"discordbot/src/lib/colours"
	"discordbot/src/lib/timeout"
)

var stoppedRegex = regexp.MustCompile(`\[.* INFO\]\: Closing Server`)

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if !timeout.GetTimeout("minecraft") {
		durationInSeconds, err := strconv.Atoi(os.Getenv("START_STOP_TIMEOUT_IN_SECONDS"))
		if err != nil {
			durationInSeconds = 30
		}
		go timeout.StartTimeout("minecraft", time.Second*time.Duration(durationInSeconds))

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
			Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Stopping...`\n`Users: None`\n",
		})

		session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      interaction.Message.ID,
			Channel: interaction.Message.ChannelID,
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: fields,
					Color:  colours.ColourOrange,
				},
			},
			Components: []discordgo.MessageComponent{},
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

		if channelId := os.Getenv("LOGS_CHANNEL_ID"); channelId != "" {
			session.ChannelMessageSendEmbed(channelId, &discordgo.MessageEmbed{
				Description: fmt.Sprintf("%s has stopped the server", interaction.Member.Mention()),
				Color:       colours.ColourRed,
			})
		}

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
