// Description: Tests for the telegramMessages package.

package telegramMessages

import (
	"testing"

	"github.com/adshao/go-binance/v2"
)

// TestNewAssetMessage tests the newAssetMessage function.
func TestNewAssetMessage(t *testing.T) {
	message := newAssetMessage("BTC", "https://www.google.com", binance.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"})
	if message != "💎 <u>Binance listed new asset (<a href='https://www.google.com'>BTC</a>)</u>\n\n- <b>Base Asset:</b> BTC\n- <b>Quota Asset:</b> USDT\n" {
		t.Errorf("Expected %s, got %s", "💎 <u>Binance listed new asset (<a href='https://www.google.com'>BTC</a>)</u>\n\n- <b>Base Asset:</b> BTC\n- <b>Quota Asset:</b> USDT\n", message)
	}
}

// TestRemovedAssetMessage tests the removedAssetMessage function.
func TestRemovedAssetMessage(t *testing.T) {
	message := removedAssetMessage("BTC")
	if message != "🗑 <u>Binance removed asset (BTC)</u>\n" {
		t.Errorf("Expected %s, got %s", "🗑 <u>Binance removed asset (BTC)</u>\n", message)
	}
}

// TestAnnouncementMessage tests the AnnouncementMessage function.
func TestAnnouncementMessage(t *testing.T) {
	message := AnnouncementMessage("https://www.google.com", "test")
	if message != "📢 <a href='https://www.google.com'>test</a>\n" {
		t.Errorf("Expected %s, got %s", "📢 <a href='https://www.google.com'>test</a>\n", message)
	}
}
