// Description: The discord package contains functions for interacting with the Discord API.
package discord

import (
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/rickstaa/crypto-listings-sniper/messaging/discord/discordEmbeds"
	"github.com/rickstaa/crypto-listings-sniper/utils"
)

// SetupDiscordSlashCommands setups the discord slash commands.
func SetupDiscordSlashCommands(discordBot *discordgo.Session, discordAppID string, telegramInviteLink string, GithubRepoURL string) {
	applicationCommands := []*discordgo.ApplicationCommand{
		{
			Name:        "telegram-invite",
			Description: "Get a invite link to the telegram channel.",
			Type:        discordgo.ChatApplicationCommand,
			Options:     []*discordgo.ApplicationCommandOption{},
		},
		{
			Name:        "github-repo",
			Description: "Get a link to the bots Github repository.",
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
		"github-repo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: GithubRepoURL,
					},
				},
			)
			if err != nil {
				log.Fatalf("Error responding to github repo slash command: %v", err)
			}
		},
	}

	// Register slash commands and handlers.
	_, err := discordBot.ApplicationCommandBulkOverwrite(discordAppID, "", applicationCommands)
	if err != nil {
		log.Fatalf("Error creating global slash commands: %v", err)
	}
	discordBot.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if i.Type == discordgo.InteractionApplicationCommand {
			name := i.ApplicationCommandData().Name
			if h, ok := commandHandlers[name]; ok {
				h(s, i)
			}
		}
	})
	discordBot.Open()
}

// sendDiscordEmbed sends a Discord embed message to the specified channel.
func sendDiscordEmbed(discordBot *discordgo.Session, discordChannelID string, embed *discordgo.MessageEmbed) {
	_, err := discordBot.ChannelMessageSendEmbed(discordChannelID, embed)
	if err != nil {
		log.Fatalf("Error sending discord embed message: %v", err)
	}
}

// SendAssetDiscordMessage sends a new/removed asset Discord embed message to a specified channel.
func SendAssetDiscordMessage(discordBot *discordgo.Session, discordChannelIDs []string, removed bool, asset string, assetInfo binance.Symbol) {
	messageEmbed := discordEmbeds.AssetEmbed(removed, asset, assetInfo)
	for _, channelID := range discordChannelIDs {
		go sendDiscordEmbed(discordBot, channelID, &messageEmbed)
	}
}

// Send Announcement Discord message to the specified channel.
func SendAnnouncementDiscordMessage(discordBot *discordgo.Session, discordChannelIDs []string, announcementCode string, announcementTitle string) {
	messageEmbed := discordEmbeds.AnnouncementEmbed(utils.CreateBinanceArticleURL(announcementCode, announcementTitle), announcementTitle)
	for _, channelID := range discordChannelIDs {
		go sendDiscordEmbed(discordBot, channelID, &messageEmbed)
	}
}
