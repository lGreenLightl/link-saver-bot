package telegram

import (
	"errors"
	"github.com/lGreenLightl/link-saver-bot/lib/e"
	"net/url"
	"strings"
)

import (
	"github.com/lGreenLightl/link-saver-bot/lib/consts"
	"github.com/lGreenLightl/link-saver-bot/storage"
)

func (p *EventsProcessor) command(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	if isAddCommand(text) {
		return p.savePage(text, chatId, username)
	}

	switch text {
	case consts.HelpCommand:
		return p.sendHelp(chatId)
	case consts.RndCommand:
		return p.sendRandom(chatId, username)
	case consts.StartCommand:
		return p.sendHello(chatId)
	default:
		return p.telegramClient.SendMessage(chatId, consts.MessageUnknownCommand)
	}
}

func (p *EventsProcessor) savePage(pageUrl string, chatId int, username string) error {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.telegramClient.SendMessage(chatId, consts.MessageAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.telegramClient.SendMessage(chatId, consts.MessageSaved); err != nil {
		return err
	}

	return nil
}

func (p *EventsProcessor) sendRandom(chatId int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return e.Wrap("can't do command: can't send random", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.telegramClient.SendMessage(chatId, consts.MessageNoSavedPages)
	}

	if err := p.telegramClient.SendMessage(chatId, page.URL); err != nil {
		return e.Wrap("can't do command: can't send random", err)
	}

	return p.storage.Remove(page)
}

func (p *EventsProcessor) sendHelp(chatId int) error {
	return p.telegramClient.SendMessage(chatId, consts.MessageHelp)
}

func (p *EventsProcessor) sendHello(chatId int) error {
	return p.telegramClient.SendMessage(chatId, consts.MessageHello)
}

func isAddCommand(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	currentUrl, err := url.Parse(text)

	return err == nil && currentUrl.Host != ""
}
