package cli

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/lewwolfe/beesinthetrap/internal/config"
	"github.com/lewwolfe/beesinthetrap/internal/game"
)

func TestPromptPlayerName(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
	}{
		{"John", "John"},
		{"", "Anonymous Hunter"},
		{"Alice", "Alice"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Set up the scanner to mock user input
			cli := NewGameCLI(nil)
			cli.scanner = bufio.NewScanner(bytes.NewReader([]byte(tt.input + "\n")))

			cli.promptPlayerName()

			// Check if the expected player name is set
			if cli.playerName != tt.expectedName {
				t.Errorf("Expected playerName to be '%s', got '%s'", tt.expectedName, cli.playerName)
			}
		})
	}
}

func TestPromptAutoMode(t *testing.T) {
	tests := []struct {
		input        string
		expectedMode bool
	}{
		{"y", true},
		{"yes", true},
		{"n", false},
		{"no", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cli := NewGameCLI(game.NewGame(config.LoadConfig()))
			cli.scanner = bufio.NewScanner(bytes.NewReader([]byte(tt.input + "\n")))

			cli.promptAutoMode()

			// Check if the expected autoMode is set
			if cli.autoMode != tt.expectedMode {
				t.Errorf("Expected autoMode to be '%v', got '%v'", tt.expectedMode, cli.autoMode)
			}
		})
	}
}
