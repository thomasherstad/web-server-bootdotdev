package main

import "testing"

func TestSilenceProfanities(t *testing.T) {
	// profanities: "kerfuffle", "sharbert", "fornax"

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no profanities",
			text:     "hello, this is houston speaking",
			expected: "hello, this is houston speaking",
		},
		{
			name:     "one profanity",
			text:     "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			name:     "two  profanities",
			text:     "This is a kerfuffle opinion I need to share with the sharbert",
			expected: "This is a **** opinion I need to share with the ****",
		},
		{
			name:     "profanity with exlamation point",
			text:     "You are a fornax!",
			expected: "You are a fornax!",
		},
		{
			name:     "capital letter",
			text:     "This is a Kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := silenceProfanities(testCase.text)
			if actual != testCase.expected {
				t.Errorf("Test %v - %s FAIL: expected: %s, actual: %s", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
