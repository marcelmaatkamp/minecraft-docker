package replies

import "github.com/bwmarrin/discordgo"

func ReplyHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func ModalReplyHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: data,
	})
}

func EditReplyHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, data *discordgo.WebhookEdit) {
	session.InteractionResponseEdit(interaction.Interaction, data)
}

func ApplyHidden(data *discordgo.InteractionResponseData, hidden bool) {
	if hidden {
		data.Flags = discordgo.MessageFlagsEphemeral
	}
}

func ModalReply(session *discordgo.Session, interaction *discordgo.InteractionCreate, applicationName string, title string, components []discordgo.MessageComponent) {
	ModalReplyHandler(session, interaction, &discordgo.InteractionResponseData{
		// Flags: discordgo.MessageFlagsEphemeral,
		CustomID:   applicationName,
		Title:      title,
		Components: components,
	})
}

func BasicThumbnailEmbedReply(session *discordgo.Session, interaction *discordgo.InteractionCreate, thumbnail string, footer string, content string, colour int, hidden bool) {
	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: thumbnail,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: footer,
				},
				Description: content,
				Color:       colour,
			},
		},
	}

	ApplyHidden(data, hidden)
	ReplyHandler(session, interaction, data)
}

func BasicEmbedReply(session *discordgo.Session, interaction *discordgo.InteractionCreate, footer string, content string, colour int, hidden bool) {
	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Footer: &discordgo.MessageEmbedFooter{
					Text: footer,
				},
				Description: content,
				Color:       colour,
			},
		},
	}

	ApplyHidden(data, hidden)
	ReplyHandler(session, interaction, data)
}

func EmbedReply(session *discordgo.Session, interaction *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, hidden bool) {
	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}

	ApplyHidden(data, hidden)
	ReplyHandler(session, interaction, data)
}

func EditBasicEmbedReply(session *discordgo.Session, interaction *discordgo.InteractionCreate, footer string, content string, colour int) {
	data := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Footer: &discordgo.MessageEmbedFooter{
					Text: footer,
				},
				Description: content,
				Color:       colour,
			},
		},
	}

	EditReplyHandler(session, interaction, data)
}
