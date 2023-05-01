package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

import (
	"github.com/lGreenLightl/link-saver-bot/lib/consts"
	"github.com/lGreenLightl/link-saver-bot/lib/e"
	"github.com/lGreenLightl/link-saver-bot/storage"
)

type FileStorage struct {
	basePath string
}

func NewStorage(basePath string) FileStorage {
	return FileStorage{basePath: basePath}
}

func (s FileStorage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		message := fmt.Sprintf("can't check if file %s exists", filePath)

		return false, e.Wrap(message, err)
	}

	return true, nil
}

func (s FileStorage) PickRandom(userName string) (*storage.Page, error) {
	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))
}

func (s FileStorage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(filePath); err != nil {
		message := fmt.Sprintf("can't remove file %s", filePath)

		return e.Wrap(message, err)
	}

	return nil
}

func (s FileStorage) Save(page *storage.Page) error {
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, consts.DefaultPermission); err != nil {
		return e.Wrap("can't save page", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return e.Wrap("can't save page", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (s FileStorage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() {
		_ = f.Close()
	}()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
