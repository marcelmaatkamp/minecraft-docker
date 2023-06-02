package start

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

var joinedRegex = regexp.MustCompile(`\[.* INFO\]\: (.*) joined the game`)
var leftRegex = regexp.MustCompile(`\[.* INFO\]\: (.*) left the game`)
var startedRegex = regexp.MustCompile(`\[.* INFO\]\: Done \(.*\)! For help, type "help"`)
var started = false
var users = []string{}

var autostopContext context.Context
var autostopCancel context.CancelFunc

func Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if !timeout.GetTimeout("minecraft") {
		durationInSeconds, err := strconv.Atoi(os.Getenv("START_STOP_TIMEOUT_IN_SECONDS"))
		if err != nil {
			durationInSeconds = 30
		}
		go timeout.StartTimeout("minecraft", time.Second*time.Duration(durationInSeconds))

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
							Value: "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Starting...`\n`Users: None`\n",
						},
					},
					Color: colours.ColourOrange,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
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
						Color:       colours.ColourBlue,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
			Description: fmt.Sprintf("%s has started the server", interaction.Member.Mention()),
			Color:       colours.ColourGreen,
		})

		users = []string{}
		started = false

		ctx, cancel := context.WithCancel(context.Background())
		autostopContext = ctx
		autostopCancel = cancel

		reader, writer := io.Pipe()

		cmdCtx, cmdDone := context.WithCancel(context.Background())

		scannerStopped := make(chan struct{})
		go func() {
			defer close(scannerStopped)

			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				line := scanner.Text()

				if !started && startedRegex.MatchString(line) {
					started = true

					message.Embeds[0].Color = colours.ColourGreen
					message.Embeds[0].Fields[3].Value = "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n`Users: None`\n"
					session.ChannelMessageEditComplex(message)

					go func(ctx context.Context) {
						durationInMinutes, err := strconv.Atoi(os.Getenv("AUTOSTOP_TIMEOUT_IN_MINUTES"))
						if err != nil {
							durationInMinutes = 30
						}

						fmt.Println("Autostop countdown starting")

						select {
						case <-ctx.Done():
							fmt.Println("Autostop countdown cancelled")
							return
						case <-time.After(time.Duration(time.Minute * time.Duration(durationInMinutes))):
							fmt.Println("Autostop initiated")

							cmd := exec.Command("pkill", "java")
							cmd.Start()

							message.Embeds[0].Color = colours.ColourRed
							message.Embeds[0].Fields[3].Value = "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Offline`\n`Users: None`\n"
							message.Components = []discordgo.MessageComponent{
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
							}

							session.ChannelMessageEditComplex(message)

							session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
								Description: "Server has stopped automatically",
								Color:       colours.ColourRed,
							})
						}
					}(autostopContext)
				}

				if joinedRegex.MatchString(line) {
					user := joinedRegex.FindStringSubmatch(line)[1]
					users = append(users, user)

					message.Embeds[0].Fields[3].Value = fmt.Sprintf("To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n`Users: %s`\n", strings.Join(users, ", "))
					session.ChannelMessageEditComplex(message)

					if autostopCancel != nil {
						autostopCancel()
					}
				}

				if leftRegex.MatchString(line) {
					user := leftRegex.FindStringSubmatch(line)[1]

					for index, searchUser := range users {
						if user == searchUser {
							users = append(users[:index], users[index+1:]...)

							if len(users) > 0 {
								message.Embeds[0].Fields[3].Value = fmt.Sprintf("To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n`Users: %s`\n", strings.Join(users, ", "))
								session.ChannelMessageEditComplex(message)

								if autostopCancel != nil {
									autostopCancel()
								}
							} else {
								message.Embeds[0].Fields[3].Value = "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Online`\n`Users: None`\n"
								session.ChannelMessageEditComplex(message)

								ctx, cancel := context.WithCancel(context.Background())
								autostopContext = ctx
								autostopCancel = cancel

								go func(ctx context.Context) {
									durationInMinutes, err := strconv.Atoi(os.Getenv("AUTOSTOP_TIMEOUT_IN_MINUTES"))
									if err != nil {
										durationInMinutes = 30
									}

									fmt.Println("Autostop countdown starting")

									select {
									case <-ctx.Done():
										fmt.Println("Autostop countdown cancelled")
										return
									case <-time.After(time.Duration(time.Minute * time.Duration(durationInMinutes))):
										fmt.Println("Autostop initiated")
										cmd := exec.Command("pkill", "java")
										cmd.Start()

										message.Embeds[0].Color = colours.ColourRed
										message.Embeds[0].Fields[3].Value = "To start/stop the minecraft server use the buttons below.\n\nWhen you want to use the server, start it and wait a minute (it boots up quickly). Once you have finished (and nobody else is using the server), please stop it.\n\n`Status: Offline`\n`Users: None`\n"
										message.Components = []discordgo.MessageComponent{
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
										}

										session.ChannelMessageEditComplex(message)

										session.ChannelMessageSendEmbed(os.Getenv("LOGS_CHANNEL"), &discordgo.MessageEmbed{
											Description: "Server has stopped automatically",
											Color:       colours.ColourRed,
										})
									}
								}(autostopContext)
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
