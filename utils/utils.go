// Description: This package contains utility functions.
package utils

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Global variables.
var (
	ASSETS_FILE_PATH = "data/assets_list.json"
)

// Delete empty strings from a slice of strings.
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// Retrieve environment variables.
func GetEnvVars() (binanceKey string, binanceSecret string, telegramBotKey string, telegramChatID int64, enableTelegramMessage bool, discordBotKey string, discordChannelIDs []string, discordAppID string, enableDiscordMessages bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	binanceKey = os.Getenv("BINANCE_API_Key")
	binanceSecret = os.Getenv("BINANCE_API_SECRET_KEY")
	telegramBotKey = os.Getenv("TELEGRAM_BOT_TOKEN")
	enableTelegramMessage, err = strconv.ParseBool(os.Getenv("ENABLE_TELEGRAM_MESSAGES"))
	if err != nil {
		log.Fatalf("Error parsing ENABLE_TELEGRAM_MESSAGES: %v", err)
	}
	discordBotKey = os.Getenv("DISCORD_BOT_TOKEN")
	telegramChatID, err = strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Error parsing TELEGRAM_CHAT_ID: %v", err)
	}
	discordChannelIDs = deleteEmpty(strings.Split(os.Getenv("DISCORD_CHANNEL_IDS"), ","))
	discordAppID = os.Getenv("DISCORD_APP_ID")
	enableDiscordMessages, err = strconv.ParseBool(os.Getenv("ENABLE_DISCORD_MESSAGES"))
	if err != nil {
		log.Fatalf("Error parsing ENABLE_DISCORD_MESSAGES: %v", err)
	}

	return binanceKey, binanceSecret, telegramBotKey, telegramChatID, enableTelegramMessage, discordBotKey, discordChannelIDs, discordAppID, enableDiscordMessages
}

// Check if a string is in a slice of strings.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Compare two lists of strings and return the difference.
func CompareLists(oldList []string, newList []string) (removed bool, difference []string) {
	if len(oldList) > len(newList) {
		for _, s := range oldList {
			if !Contains(newList, s) {
				difference = append(difference, s)
			}
		}

		return true, difference
	} else if len(oldList) < len(newList) {
		// If the length of the new list is greater than the old list.
		for _, s := range newList {
			removed = false
			if !Contains(oldList, s) {
				difference = append(difference, s)
			}
		}

		return false, difference
	}

	return false, []string{}
}

// Retrieve the old assets from the data folder.
func RetrieveOldListings() (oldAssets []string) {
	oldAssetsJson, err := os.ReadFile(ASSETS_FILE_PATH)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading old assets from file '%s': %v", ASSETS_FILE_PATH, err)
		}
	} else {
		err = json.Unmarshal(oldAssetsJson, &oldAssets)
		if err != nil {
			log.Fatalf("Error unmarshalling old assets: %v", err)
		}
		log.Printf("Number of old assets: %d", len(oldAssets))
	}

	return oldAssets
}

// Store the old assets in the data folder.
func StoreOldListings(assets []string) {
	if len(assets) != 0 {
		assetsJson, err := json.Marshal(assets)
		if err != nil {
			log.Fatalf("Error marshalling assets list: %v", err)
		}
		err = os.WriteFile(ASSETS_FILE_PATH, assetsJson, 0644)
		if err != nil {
			log.Fatalf("Error writing assets list to file '%s': %v", ASSETS_FILE_PATH, err)
		}
	}
}

// Create Binance URL for a asset.
func CreateBinanceURL(assetName string) string {
	return "https://www.binance.com/en/trade/" + assetName
}
