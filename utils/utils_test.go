// Description: Tests for the utils package.

package utils

import (
	"testing"
)

// TestGetEnvVar tests the GetEnvVar function.
func TestDeleteEmpty(t *testing.T) {
	s := []string{"", "hello", "", "world", ""}
	r := deleteEmpty(s)
	if len(r) != 2 {
		t.Errorf("Expected length of 2, got %d", len(r))
	}
}

// TestContains tests the Contains function.
func TestContains(t *testing.T) {
	s := []string{"hello", "world"}
	if !contains(s, "hello") {
		t.Errorf("Expected true, got false")
	}
	if contains(s, "test") {
		t.Errorf("Expected false, got true")
	}
}

// TestHexColorToInt tests the HexColorToInt function.
func TestHexColorToInt(t *testing.T) {
	color := "#ff0000"
	expected := 16711680
	r := HexColorToInt(color)
	if r != expected {
		t.Errorf("Expected %d, got %d", expected, r)
	}

	color = "ff0000"
	expected = 16711680
	r = HexColorToInt(color)
	if r != expected {
		t.Errorf("Expected %d, got %d", expected, r)
	}
}

// TestCompareLists tests the CompareLists function.
func TestCompareLists(t *testing.T) {
	list1 := []string{"hello", "world"}
	list2 := []string{"hello", "world"}
	r, _ := CompareLists(list1, list2)
	if r {
		t.Errorf("Expected false, got true")
	}

	list1 = []string{"hello", "world"}
	list2 = []string{"hello", "world", "test"}
	r, diff := CompareLists(list1, list2)
	if r {
		t.Errorf("Expected false, got true")
	}
	if len(diff) != 1 {
		t.Errorf("Expected length of 1, got %d", len(diff))
	}
	list1 = []string{"hello", "world", "test"}
	list2 = []string{"hello", "world"}
	r, diff = CompareLists(list1, list2)
	if !r {
		t.Errorf("Expected true, got false")
	}
	if len(diff) != 1 {
		t.Errorf("Expected length of 1, got %d", len(diff))
	}
}

// TestCreateBinanceURL tests the CreateBinanceURL function.
func TestCreateBinanceURL(t *testing.T) {
	expected := "https://www.binance.com/en/trade/BLC"
	r := CreateBinanceURL("BLC")
	if r != expected {
		t.Errorf("Expected %s, got %s", expected, r)
	}
}

// TestCreateBinanceArticleUrl tests the CreateBinanceArticleUrl function.
func TestCreateBinanceArticleUrl(t *testing.T) {
	expected := "https://www.binance.com/en/support/announcement/article-48"
	r := CreateBinanceArticleURL("48", "article")
	if r != expected {
		t.Errorf("Expected %s, got %s", expected, r)
	}
}
