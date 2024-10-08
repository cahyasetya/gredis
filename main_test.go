package main

import (
	"encoding/binary"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "Single command",
			input:    createMessage([]string{"SET", "key", "value"}),
			expected: []string{"SET", "key", "value"},
		},
		{
			name:     "Multiple commands",
			input:    createMessage([]string{"GET", "key", "DEL", "key", "PING"}),
			expected: []string{"GET", "key", "DEL", "key", "PING"},
		},
		{
			name:     "Empty command",
			input:    createMessage([]string{""}),
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCommand(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d commands, got %d", len(tt.expected), len(result))
			}

			for i, cmd := range result {
				if i >= len(tt.expected) {
					break
				}
				if cmd != tt.expected[i] {
					t.Errorf("Expected command %q, got %q", tt.expected[i], cmd)
				}
			}
		})
	}
}

func createMessage(commands []string) []byte {
	var result []byte
	cmdCount := uint32(len(commands))
	countBytes := make([]byte, LenSize)
	binary.LittleEndian.PutUint32(countBytes, cmdCount)
	result = append(result, countBytes...)

	for _, cmd := range commands {
		length := uint32(len(cmd))
		lengthBytes := make([]byte, LenSize)
		binary.LittleEndian.PutUint32(lengthBytes, length)
		result = append(result, lengthBytes...)
		result = append(result, []byte(cmd)...)
	}
	return result
}
