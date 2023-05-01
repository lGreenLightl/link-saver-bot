package main

import (
	"flag"
	"log"
)

func main() {}

func mustToken() string {
	token := flag.String(
		"telegram-bot-token",
		"",
		"telegram access token")
	flag.Parse()

	if *token == "" {
		log.Fatal("token isn't specified")
	}

	return *token
}
