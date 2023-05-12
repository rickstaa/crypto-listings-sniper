// Description: Contains functions for interacting with the Discord API.
package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rickstaa/crypto-listings-sniper/discord/discordEmbeds"
)

// Setup discord slash commands.
func SetupDiscordSlashCommands(discordBot *discordgo.Session, discordAppID string, telegramInviteLink string) {
	applicationCommands := []*discordgo.ApplicationCommand{
		{
			Name:        "telegram-invite",
			Description: "Get a invite link to the telegram channel.",
			Type:        discordgo.ChatApplicationCommand,
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"telegram-invite": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: telegramInviteLink,
					},
				},
			)
			if err != nil {
				log.Fatalf("Error responding to telegram invite slash command: %v", err)
			}
		},
	}

	// Register slash commands and handlers.
	_, err := discordBot.ApplicationCommandBulkOverwrite(discordAppID, "", applicationCommands)
	if err != nil {
		log.Fatalf("Error creating invite slash command: %v", err)
	}
	discordBot.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			name := i.ApplicationCommandData().Name
			if h, ok := commandHandlers[name]; ok {
				h(s, i)
			}
		}
	})
	discordBot.Open()
}

// Send a Discord embed message to the specified channel.
func sendDiscordEmbed(discordBot *discordgo.Session, discordChannelID string, embed *discordgo.MessageEmbed) {
	_, err := discordBot.ChannelMessageSendEmbed(discordChannelID, embed)
	if err != nil {
		log.Fatalf("Error sending discord embed message: %v", err)
	}
}

// Send Trading pair Discord message to the specified channel.
func SendAssetDiscordMessage(discordBot *discordgo.Session, discordChannelIDs []string, removed bool, asset string) {
	messageEmbed := discordEmbeds.AssetEmbed(removed, asset)
	for _, channelID := range discordChannelIDs {
		go sendDiscordEmbed(discordBot, channelID, &messageEmbed)
	}
}
