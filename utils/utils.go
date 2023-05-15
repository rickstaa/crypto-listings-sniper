// Description: The utils package contains contains utility functions.
package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Global variables.
var (
	ASSETS_FILE_PATH        = "data/assets_list.json"
	ANNOUNCEMENTS_FILE_PATH = "data/announcements_list.json"
)

// deleteEmpty deletes empty strings from a slice of strings.
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}

// GetEnvVars retrieves environment variables.
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

// Contains checks if a string is in a slice of strings.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// CompareLists compares two lists of strings and returns the items that were removed/added in the new list.
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

// RetrieveOldListings retrieves the old listed assets from the data folder.
func RetrieveOldListings() (oldAssets []string) {
	oldAssetsJson, err := os.ReadFile(ASSETS_FILE_PATH)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading old listed assets from file '%s': %v", ASSETS_FILE_PATH, err)
		}
	} else {
		err = json.Unmarshal(oldAssetsJson, &oldAssets)
		if err != nil {
			log.Fatalf("Error unmarshalling old listed assets: %v", err)
		}
		log.Printf("Number of old listed assets: %d", len(oldAssets))
	}

	return oldAssets
}

// StoreOldListings stores the old listed assets in the data folder.
func StoreOldListings(assets []string) {
	if len(assets) != 0 {
		assetsJson, err := json.Marshal(assets)
		if err != nil {
			log.Fatalf("Error marshalling listed assets: %v", err)
		}
		err = os.WriteFile(ASSETS_FILE_PATH, assetsJson, 0644)
		if err != nil {
			log.Fatalf("Error writing listed assets list to file '%s': %v", ASSETS_FILE_PATH, err)
		}
	}
}

// RetrieveOldAnnouncements retrieves the old announcements from the data folder.
func RetrieveOldAnnouncements() (oldAnnouncements []string) {
	oldAssetsJson, err := os.ReadFile(ANNOUNCEMENTS_FILE_PATH)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading old announcements from file '%s': %v", ANNOUNCEMENTS_FILE_PATH, err)
		}
	} else {
		err = json.Unmarshal(oldAssetsJson, &oldAnnouncements)
		if err != nil {
			log.Fatalf("Error unmarshalling old announcements: %v", err)
		}
		log.Printf("Number of old assets: %d", len(oldAnnouncements))
	}

	return oldAnnouncements
}

// StoreOldAnnouncements store the old assets in the data folder.
func StoreOldAnnouncements(announcements []string) {
	if len(announcements) != 0 {
		announcementsJson, err := json.Marshal(announcements)
		if err != nil {
			log.Fatalf("Error marshalling announcements list: %v", err)
		}
		err = os.WriteFile(ANNOUNCEMENTS_FILE_PATH, announcementsJson, 0644)
		if err != nil {
			log.Fatalf("Error writing announcements list to file '%s': %v", ANNOUNCEMENTS_FILE_PATH, err)
		}
	}
}

// Create Binance URL for a asset.
func CreateBinanceURL(assetName string) string {
	return "https://www.binance.com/en/trade/" + assetName
}

// CreateBinanceArticleURL returns the binance article URL.
func CreateBinanceArticleURL(articleCode string, articleTitle string) string {
	// Replace spaces with hypens and make lowercase.
	articleTitle = strings.ReplaceAll(articleTitle, " ", "-")
	articleTitle = strings.ToLower(articleTitle)

	return fmt.Sprintf("https://www.binance.com/en/support/announcement/%s-%s", articleTitle, articleCode)
}

// GetBinanceAnnouncementsEndpoint returns the (unofficial) binance announcements endpoint.
func GetBinanceAnnouncementsEndpoint() string {
	queries := map[string]string{
		"catalogId": "48",
		"pageNo":    "1",
		"pageSize":  "10", // NOTE: This is to prevent the endpoint from being cached.
	}
	var url strings.Builder
	url.WriteString("https://www.binance.com/bapi/composite/v1/public/cms/article/catalog/list/query?")
	for key, value := range queries {
		url.WriteString(key + "=" + value + "&")
	}
	return url.String()
}

// HexColorToInt converts a hex color to int.
func HexColorToInt(color string) int {
	colorInt, err := strconv.ParseUint(color, 16, 64)
	if err != nil {
		log.Fatalf("Error parsing color: %v", err)
	}
	return int(colorInt)
}
