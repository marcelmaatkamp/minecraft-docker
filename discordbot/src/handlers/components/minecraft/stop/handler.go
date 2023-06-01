package stop

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bwmarrin/discordgo"

	dsUtils "discordbot/src/lib/utils"
)

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
						Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Offline`\n\n`Users: None`\n",
					},
				},
				Color: dsUtils.ColourRed,
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
					Color:       dsUtils.ColourBlue,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s has stopped the server", interaction.Member.Mention()),
		Color:       dsUtils.ColourRed,
	})

	reader, writer := io.Pipe()

	cmdCtx, cmdDone := context.WithCancel(context.Background())

	scannerStopped := make(chan struct{})
	go func() {
		defer close(scannerStopped)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	cmd := exec.Command("pkill", "java")
	cmd.Stdout = writer
	_ = cmd.Start()
	go func() {
		_ = cmd.Wait()
		cmdDone()
		writer.Close()
	}()
	<-cmdCtx.Done()

	<-scannerStopped
}
