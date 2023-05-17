package main

import (
	"log"
	"runtime"

	"github.com/bwmarrin/discordgo"

	"github.com/rickstaa/crypto-listings-sniper/exchanges/binanceAnnouncementsChecker"
	"github.com/rickstaa/crypto-listings-sniper/exchanges/binanceListingsChecker"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"
	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	envVars := utils.GetEnvVars()

	// Load Telegram bot.
	telegramBot, err := telego.NewBot(envVars.TelegramBotKey)
	if err != nil {
		log.Fatalf("Error loading Telegram telegramBot: %v", err)
	}

	// Load Discord bot.
	discordBot, err := discordgo.New("Bot " + envVars.DiscordBotKey)
	if err != nil {
		log.Fatalf("Error loading Discord bot: %v", err)
	}

	// Log Telegram bot and channel info.
	telegramBotInfo, err := telegramBot.GetMe()
	if err != nil {
		log.Fatalf("Error getting telegramBot info: %v", err)
	}
	telegramChat, err := telegramBot.GetChat(&telego.GetChatParams{ChatID: tu.ID(envVars.TelegramChatID)})
	if err != nil {
		log.Fatalf("Error getting telegramChat info: %v", err)
	}
	log.Printf("Authorized on account: %s", telegramBotInfo.Username)
	log.Printf("Bot id: %d", telegramBotInfo.ID)
	log.Printf("Chat id: %d", envVars.TelegramChatID)
	log.Printf("Chat type: %s", telegramChat.Type)
	log.Printf("Chat title: %s", telegramChat.Title)
	log.Printf("Chat username: %s", telegramChat.Username)
	log.Printf("Chat description: %s", telegramChat.Description)

	// Register slash commands.
	dc.SetupDiscordSlashCommands(discordBot, envVars.DiscordAppID, telegramChat.InviteLink)

	// Initialize Binance client.
	binanceClient := binance.NewClient(envVars.BinanceKey, envVars.BinanceSecret)
	binanceClient.SetApiEndpoint("https://api4.binance.com")
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)
	log.Printf("Binance announcement API endpoint: %s", binanceAnnouncementsChecker.GetBinanceAnnouncementsEndpoint())

	// Initialize crypto checkers.
	binanceListingsChecker := binanceListingsChecker.NewBinanceListingsChecker(binanceClient, telegramBot, envVars.TelegramChatID, envVars.EnableTelegramMessage, discordBot, envVars.DiscordChannelIDs, envVars.EnableDiscordMessages)
	binanceAnnouncementsChecker := binanceAnnouncementsChecker.NewBinanceAnnouncementsChecker(binanceClient, telegramBot, envVars.TelegramChatID, envVars.EnableTelegramMessage, discordBot, envVars.DiscordChannelIDs, envVars.EnableDiscordMessages)

	// start the checkers.
	go binanceListingsChecker.Start(envVars.BinanceListingsRate)
	go binanceAnnouncementsChecker.Start(envVars.BinanceAnnouncementsRate)

	runtime.Goexit() // Keep the program running.
}
