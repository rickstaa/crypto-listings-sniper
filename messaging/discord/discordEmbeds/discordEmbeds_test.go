// Description: Tests for the discrodEmbeds package.

package discordEmbeds

import (
	"testing"

	"github.com/adshao/go-binance/v2"
)

// TestNewAssetMessage tests the newAssetMessage function.
func TestNewAssetMessage(t *testing.T) {
	embed := newAssetMessage("BTC", binance.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"})
	if embed.Title != "💎 Binance listed new asset (BTC)" {
		t.Errorf("Expected %s, got %s", "💎 Binance listed new asset (BTC)", embed.Title)
	}
	if embed.Description != "• **Base Asset:** BTC\n• **Quota Asset:** USDT\n" {
		t.Errorf("Expected %s, got %s", "• **Base Asset:** BTC\n• **Quota Asset:** USDT\n", embed.Description)
	}
}

// TestRemovedAssetMessage tests the removedAssetMessage function.
func TestRemovedAssetMessage(t *testing.T) {
	embed := removedAssetMessage("BTC")
	if embed.Title != "🗑 Binance removed asset (BTC)\n" {
		t.Errorf("Expected %s, got %s", "🗑 Binance removed asset (BTC)\n", embed.Title)
	}
	if embed.Image != nil {
		t.Errorf("Expected %v, got %v", nil, embed.Image)
	}
}

// TestAssetEmbed tests the AssetEmbed function.
func TestAnnouncementEmbed(t *testing.T) {
	embed := AnnouncementEmbed("https://www.google.com", "Test")
	if embed.Title != "📢 Test" {
		t.Errorf("Expected %s, got %s", "📢 Test", embed.Title)
	}
	if embed.URL != "https://www.google.com" {
		t.Errorf("Expected %s, got %s", "https://www.google.com", embed.URL)
	}
}
