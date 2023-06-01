package utils

import (
	"github.com/bwmarrin/discordgo"
)

func GetUser(interaction *discordgo.InteractionCreate) (*discordgo.User, bool) {
	dm := true
	user := interaction.User
	if user == nil {
		user = interaction.Member.User
		dm = false
	}
	return user, dm
}

func InArray(array []any, value any) bool {
	for child := range array {
		if child == value {
			return true
		}
	}
	return false
}
