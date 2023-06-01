package start

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	dsUtils "discordbot/src/lib/utils"
)

var joinedRegex = regexp.MustCompile(`\[.* INFO\]\: (.*) joined the game`)
var leftRegex = regexp.MustCompile(`\[.* INFO\]\: (.*) left the game`)
var users = []string{}

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	message := &discordgo.MessageEdit{
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
						Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n\n`Users: None`\n",
					},
				},
				Color: dsUtils.ColourGreen,
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
	}

	session.ChannelMessageEditComplex(message)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: "Starting minecraft server...",
					Color:       dsUtils.ColourBlue,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s has started the server", interaction.Member.Mention()),
		Color:       dsUtils.ColourGreen,
	})

	reader, writer := io.Pipe()

	cmdCtx, cmdDone := context.WithCancel(context.Background())

	scannerStopped := make(chan struct{})
	go func() {
		defer close(scannerStopped)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			if joinedRegex.MatchString(line) {
				user := joinedRegex.FindStringSubmatch(line)[1]
				users = append(users, user)

				message.Embeds[0].Fields[3].Value = fmt.Sprintf("To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n\n`Users: %s`\n", strings.Join(users, ", "))
				session.ChannelMessageEditComplex(message)
			}

			if leftRegex.MatchString(line) {
				user := leftRegex.FindStringSubmatch(line)[1]

				for index, searchUser := range users {
					if user == searchUser {
						users = append(users[:index], users[index+1:]...)

						if len(users) > 0 {
							message.Embeds[0].Fields[3].Value = fmt.Sprintf("To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n\n`Users: %s`\n", strings.Join(users, ", "))
							session.ChannelMessageEditComplex(message)
						} else {
							message.Embeds[0].Fields[3].Value = "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n\n`Users: None`\n"
							session.ChannelMessageEditComplex(message)
						}

						break
					}
				}
			}

			fmt.Println(line)
		}
	}()

	cmd := exec.Command("/bin/bash", "/scripts/start_java.sh")
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
