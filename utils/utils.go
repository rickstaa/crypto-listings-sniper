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
	SYMBOLS_FILE_PATH = "data/symbol_list.json"
	ASSETS_FILE_PATH  = "data/assets_list.json"
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
func GetEnvVars() (telegramBotKey string, chatID int64, binanceKey string, binanceSecret string, discordBotKey string, discordChannelIDs []string, discordAppID string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	telegramBotKey = os.Getenv("TELEGRAM_BOT_TOKEN")
	binanceKey = os.Getenv("BINANCE_API_Key")
	binanceSecret = os.Getenv("BINANCE_API_SECRET_KEY")
	discordBotKey = os.Getenv("DISCORD_BOT_TOKEN")
	chatID, err = strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Error parsing TELEGRAM_CHAT_ID: %v", err)
	}
	discordChannelIDs = deleteEmpty(strings.Split(os.Getenv("DISCORD_CHANNEL_IDS"), ","))
	discordAppID = os.Getenv("DISCORD_APP_ID")

	return telegramBotKey, chatID, binanceKey, binanceSecret, discordBotKey, discordChannelIDs, discordAppID
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

// Retrieve added or removed assets.
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

// Retrieve the old assets and symbols from the data folder.
func RetrieveOldListings() (oldBaseAssetsList []string, oldSymbolList []string) {
	oldBaseAssetsListJson, err := os.ReadFile(ASSETS_FILE_PATH)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading old 'SPOT' base assets from file '%s': %v", ASSETS_FILE_PATH, err)
		}
	} else {
		err = json.Unmarshal(oldBaseAssetsListJson, &oldBaseAssetsList)
		if err != nil {
			log.Fatalf("Error unmarshalling old base assets: %v", err)
		}
		log.Printf("Number of old 'SPOT' base assets: %d", len(oldBaseAssetsList))
	}
	oldSymbolListJson, err := os.ReadFile(SYMBOLS_FILE_PATH)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading old 'SPOT' symbols from file '%s': %v", SYMBOLS_FILE_PATH, err)
		}
	} else {
		err = json.Unmarshal(oldSymbolListJson, &oldSymbolList)
		if err != nil {
			log.Fatalf("Error unmarshalling old symbols: %v", err)
		}
		log.Printf("Number of old 'SPOT' symbols: %d", len(oldSymbolList))
	}

	return oldBaseAssetsList, oldSymbolList
}

// Store the old assets and symbols from the data folder.
func StoreOldListings(baseAssets []string, symbolList []string) {
	// Store SPOT base assets.
	if len(baseAssets) != 0 {
		baseAssetsListJson, err := json.Marshal(baseAssets)
		if err != nil {
			log.Fatalf("Error marshalling base assets list: %v", err)
		}
		err = os.WriteFile(ASSETS_FILE_PATH, baseAssetsListJson, 0644)
		if err != nil {
			log.Fatalf("Error writing base assets list to file '%s': %v", ASSETS_FILE_PATH, err)
		}
	}

	// Store SPOT symbols.
	if len(symbolList) != 0 {
		symbolListJson, err := json.Marshal(symbolList)
		if err != nil {
			log.Fatalf("Error marshalling symbol list: %v", err)
		}
		err = os.WriteFile(SYMBOLS_FILE_PATH, symbolListJson, 0644)
		if err != nil {
			log.Fatalf("Error writing symbol list to file '%s': %v", SYMBOLS_FILE_PATH, err)
		}
	}
}

// Create Binance URL for a asset.
func CreateBinanceURL(assetName string) string {
	return "https://www.binance.com/en/trade/" + assetName
}
