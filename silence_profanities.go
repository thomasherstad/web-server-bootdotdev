package main

import (
	"slices"
	"strings"
)

func silenceProfanities(content string) string {
	words := strings.Fields(content)

	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for i, word := range words {
		if slices.Contains(profanities, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
