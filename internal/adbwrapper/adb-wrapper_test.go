package adbwrapper

import (
	"testing"
)

func TestDevicesStdOutToDevices(t *testing.T) {
	result := DevicesStdOutToDevices([]byte("List of devices attached\nRFCW32969NW     device\nRFCW32969NE     device\n\n"))
	if result == nil {
		t.Errorf("Expected nil, got %v", result)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2, got %v", len(result))
	}

	if result[0].Name != "RFCW32969NW" {
		t.Errorf("Expected device 1 to have device ID RFCW32969NW, got %v", result[0].Name)
	}

	if result[1].Name != "RFCW32969NE" {
		t.Errorf("Expected device 2 to have device ID RFCW32969NE, got %v", result[1].Name)
	}
}
