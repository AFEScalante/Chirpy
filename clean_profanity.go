package main

import "strings"

func cleanProfanity(input string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, badWord := range profaneWords {
		words := strings.Fields(input)
		for i, w := range words {
			if strings.ToLower(w) == badWord {
				words[i] = "****"
			}
		}
		input = strings.Join(words, " ")
	}

	return input
}
