package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"storage-links-bot/lib/e"
	"storage-links-bot/storage"
	"time"
)

const (
	saveErr     = "cannot save file"
	checkErr    = "can't check file: %s"
	openFileErr = "cannot open file"
	openDirErr  = "cannot open directory"
	searchErr   = "not found files"
	removeErr   = "cannot remove file: %s"
	decodeErr   = "cannot decode file: %s"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath}
}

func (s Storage) Save(p *storage.Page) error {
	fPath := filepath.Join(s.basePath, p.UserName)
	if err := os.MkdirAll(fPath, 0755); err != nil {
		return e.Wrap(saveErr, err)
	}
	fName, err := p.Hash()
	if err != nil {
		return err
	}
	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap(saveErr, err)
	}
	defer file.Close()
	err = gob.NewDecoder(file).Decode(p)
	if err != nil {
		return e.Wrap(saveErr, err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap(openDirErr, err)
	}
	if len(dir) == 0 {
		return nil, errors.New(searchErr)
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(dir))
	file := dir[n]
	return decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := p.Hash()
	if err != nil {
		return err
	}
	fullName := filepath.Join(s.basePath, p.UserName, fileName)
	err = os.Remove(fullName)
	if err != nil {
		return e.Wrap(fmt.Sprintf(removeErr, fullName), err)
	}
	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := p.Hash()
	if err != nil {
		return false, err
	}
	fullName := filepath.Join(s.basePath, p.UserName, fileName)
	_, err = os.Stat(fullName)
	switch _, err = os.Stat(fullName); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf(checkErr, fullName)
		return false, e.Wrap(msg, err)
	default:
		return true, nil
	}
}

func decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(openFileErr, err)
	}
	defer file.Close()
	var p storage.Page
	err = gob.NewEncoder(file).Encode(&p)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf(decodeErr, filePath), err)
	}
	return &p, nil
}
