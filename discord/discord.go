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
}

// Send Base asset Discord message to the specified channel.
func SendBaseAssetDiscordMessage(discordBot *discordgo.Session, discordChannelIDs []string, removed bool, symbol string) {
	messageEmbed := discordEmbeds.BaseAssetEmbed(removed, symbol)

	for _, channelID := range discordChannelIDs {
		go discordBot.ChannelMessageSendEmbed(channelID, &messageEmbed)
	}
}

// Send Trading pair Discord message to the specified channel.
func SendTradingPairDiscordMessage(discordBot *discordgo.Session, discordChannelIDs []string, removed bool, symbol string) {
	messageEmbed := discordEmbeds.TradingPairEmbed(removed, symbol)
	for _, channelID := range discordChannelIDs {
		go discordBot.ChannelMessageSendEmbed(channelID, &messageEmbed)
	}
}

// Handle discord '/invite' slash command.
func HandleDiscordInviteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create the interaction response.
	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "https://discord.com/api/oauth2/authorize?client_id=889267457513744680&permissions=0&scope=bot%20applications.commands",
		},
	}
	// Send the interaction response.
	s.InteractionRespond(i.Interaction, &response)
}