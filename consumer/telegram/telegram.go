package telegram

import (
	"log"
	"time"
)

import "github.com/lGreenLightl/link-saver-bot/events"

type EventsConsumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewEventsConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) EventsConsumer {
	return EventsConsumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *EventsConsumer) Start() error {
	for {
		curEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(curEvents) == 0 {
			time.Sleep(time.Second * 1)

			continue
		}

		if err := c.handleEvents(curEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *EventsConsumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
