// Description: This package contains utility functions.
package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/rickstaa/crypto-listings-sniper/utils/telegramMessages"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Global variables.
var (
	SYMBOLS_FILE_PATH = "data/symbol_list.json"
	ASSETS_FILE_PATH  = "data/assets_list.json"
)

// Retrieve environment variables.
func GetEnvVars() (telegramBotKey string, chatID int64, binanceKey string, binanceSecret string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	telegramBotKey = os.Getenv("TELEGRAM_BOT_TOKEN")
	binanceKey = os.Getenv("BINANCE_API_Key")
	binanceSecret = os.Getenv("BINANCE_API_SECRET_KEY")
	chatID, err = strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Error parsing TELEGRAM_CHAT_ID: %v", err)
	}

	return telegramBotKey, chatID, binanceKey, binanceSecret
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
	removed = false
	var compareBase = newList
	var compareChild = oldList
	if len(oldList) > len(newList) {
		removed = true
		compareBase = oldList
		compareChild = newList
	}

	for _, s := range compareBase {
		if !Contains(compareChild, s) {
			difference = append(difference, s)
		}
	}

	return removed, difference
}

// Send a Telegram message to the specified chat.
func SendTelegramMessage(telegramBot *telego.Bot, chatID int64, message string) {
	msg := tu.Message(tu.ID(chatID), message)
	msg.ParseMode = telego.ModeHTML
	_, err := telegramBot.SendMessage(msg)
	if err != nil {
		log.Printf("WARNING: Error sending message '%s' to channel '%d': %v", msg.Text, chatID, err)
	}
}

// Send a base asset Telegram message to the specified chat.
func SendBaseAssetTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, symbol string) {
	message := telegramMessages.BaseAssetMessage(removed, symbol, createBinanceURL(symbol)+"_USDT")
	SendTelegramMessage(telegramBot, chatID, message)
}

// Send a trading pair Telegram message to the specified chat.
func SendTradingPairTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, symbol string, symbolInfo map[string]binance.Symbol) {
	message := telegramMessages.TradingPairMessage(removed, createBinanceURL(symbol), symbolInfo[symbol].BaseAsset+"/"+symbolInfo[symbol].QuoteAsset)
	SendTelegramMessage(telegramBot, chatID, message)
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

// Retrieve SPOT assets and symbols from Binance.
func RetrieveBinanceSpotAssets(binanceClient *binance.Client) (baseAssets []string, symbols []string, symbolInfo map[string]binance.Symbol) {
	// Retrieve SPOT exchange info.
	exchangeInfoService := binanceClient.NewExchangeInfoService()
	exchangeInfoService = exchangeInfoService.Permissions("SPOT")
	exchangeInfo, err := exchangeInfoService.Do(context.Background())

	// Retrieve SPOT symbols.
	exchangeSymbols := exchangeInfo.Symbols
	if err != nil {
		log.Fatalf("Error retrieving Binance 'SPOT' symbols: %v", err)
	}
	symbols = make([]string, len(exchangeSymbols))
	symbolInfo = map[string]binance.Symbol{}
	for i, s := range exchangeSymbols {
		symbols[i] = s.Symbol
		symbolInfo[s.Symbol] = s
	}

	// Retrieve SPOT base assets.
	k := make(map[string]bool)
	for _, s := range exchangeSymbols {
		if _, value := k[s.BaseAsset]; !value {
			k[s.BaseAsset] = true
			baseAssets = append(baseAssets, s.BaseAsset)
		}
	}

	return baseAssets, symbols, symbolInfo
}

// Create Binance URL for a asset.
func createBinanceURL(assetName string) string {
	return "https://www.binance.com/en/trade/" + assetName
}
