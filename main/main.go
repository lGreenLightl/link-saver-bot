package main

import (
	"flag"
	"log"
)

import (
	tgClient "github.com/lGreenLightl/link-saver-bot/clients/telegram"
	tgEventsConsumer "github.com/lGreenLightl/link-saver-bot/consumer/telegram"
	tgEventsProcessor "github.com/lGreenLightl/link-saver-bot/events/telegram"
	"github.com/lGreenLightl/link-saver-bot/lib/consts"
	"github.com/lGreenLightl/link-saver-bot/storage/files"
)

func main() {
	telegramClient := tgClient.NewClient(consts.HostPath, mustToken())

	eventsProcessor := tgEventsProcessor.NewEventsProcessor(telegramClient, files.NewStorage(consts.StoragePath))

	log.Print("service started")

	eventsConsumer := tgEventsConsumer.NewEventsConsumer(eventsProcessor, eventsProcessor, consts.BatchSize)
	if err := eventsConsumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

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
