package telegram

import (
	"errors"
	"github.com/lGreenLightl/link-saver-bot/clients/telegram"
	"github.com/lGreenLightl/link-saver-bot/events"
	"github.com/lGreenLightl/link-saver-bot/lib/consts"
	"github.com/lGreenLightl/link-saver-bot/lib/e"
	"github.com/lGreenLightl/link-saver-bot/storage"
)

type EventsProcessor struct {
	offset         int
	telegramClient *telegram.Client
	storage        storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewEventProcessor(telegramClient *telegram.Client, storage storage.Storage) *EventsProcessor {
	return &EventsProcessor{
		telegramClient: telegramClient,
		storage:        storage,
	}
}

func (p *EventsProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.telegramClient.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		result = append(result, toEvent(update))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return result, nil
}

func (p *EventsProcessor) Process(event events.Event) error {
	switch event.Type {
	case consts.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *EventsProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.command(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return result, nil
}

func toEvent(update telegram.Update) events.Event {
	updateType := fetchType(update)

	result := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == consts.Message {
		result.Meta = Meta{
			ChatId:   update.Message.Chat.Id,
			Username: update.Message.From.Username,
		}
	}

	return result
}

func fetchType(update telegram.Update) consts.Type {
	if update.Message == nil {
		return consts.Unknown
	}

	return consts.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
