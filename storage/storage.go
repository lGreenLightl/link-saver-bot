package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

import (
	"github.com/lGreenLightl/link-saver-bot/lib/e"
)

type Storage interface {
	IsExists(page *Page) (bool, error)
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	Save(page *Page) error
}

type Page struct {
	URL      string
	UserName string
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	pageHash := sha1.New()

	if _, err := io.WriteString(pageHash, p.URL); err != nil {
		return "", e.Wrap("can't calculate page hash", err)
	}

	if _, err := io.WriteString(pageHash, p.UserName); err != nil {
		return "", e.Wrap("can't calculate page hash", err)
	}

	return fmt.Sprintf("%x", pageHash.Sum(nil)), nil
}
