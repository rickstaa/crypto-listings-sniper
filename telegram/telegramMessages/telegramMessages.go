// Description: This package contains functions for creating telegram messages.
package telegramMessages

import (
	"fmt"
)

// Returns a string containing a message for a new SPOT trading pair.
func newTradingPairMessage(symbol string, url string) string {
	return fmt.Sprintf("⚖️ <u>Binance listed new SPOT trading pair (<a href='%s'>%s</a>)</u>", url, symbol)
}

// Return a string containing a message for a removed SPOT trading pair.
func removedTradingPairMessage(symbol string) string {
	return fmt.Sprintf("🗑 <u>Binance removed SPOT trading pair (%s)</u>\n", symbol)
}

// Returns a string containing a message for a new SPOT base asset.
func newBaseAssetMessage(symbol string, url string) string {
	return fmt.Sprintf("💎 <u>Binance listed new SPOT asset (<a href='%s'>%s</a>)</u>", url, symbol)
}

// Returns a string containing a message for a removed SPOT base asset.
func removedBaseAssetMessage(symbol string) string {
	return fmt.Sprintf("🗑 <u>Binance removed SPOT asset (%s)</u>\n", symbol)
}

// Returns a string containing a telegram message for a new or removed SPOT base asset.
func BaseAssetMessage(removed bool, symbol string, url string) string {
	switch removed {
	case true:
		return removedBaseAssetMessage(symbol)
	default:
		return newBaseAssetMessage(symbol, url)
	}
}

// Returns a string containing a telegram message for a new or removed SPOT trading pair.
func TradingPairMessage(removed bool, symbol string, url string) string {
	switch removed {
	case true:
		return removedTradingPairMessage(symbol)
	default:
		return newTradingPairMessage(symbol, url)
	}
}
