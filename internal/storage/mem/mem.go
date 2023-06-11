package mem

import (
	"github.com/IKostarev/yandex-go-dev/internal/model"
	"github.com/IKostarev/yandex-go-dev/internal/utils"
	uuID "github.com/google/uuid"
)

type Mem struct {
	cacheMemory      map[string]string
	cacheCorrelation map[string]string
}

func NewMem() (*Mem, error) {
	m := &Mem{
		cacheMemory:      make(map[string]string),
		cacheCorrelation: make(map[string]string),
	}

	return m, nil
}

func (m *Mem) Save(long, corrID string, _ uuID.UUID) (string, error) {
	short := utils.RandomString()

	m.cacheMemory[short] = long
	m.cacheCorrelation[corrID] = long

	return short, nil
}

func (m *Mem) Get(short, corrID string) (string, string) {
	return m.cacheMemory[short], corrID
}

func (m *Mem) GetUserLinks(_ uuID.UUID) (data []model.UserLink, err error) {
	data = make([]model.UserLink, 0)

	for originalURL, shortURL := range m.cacheMemory {
		data = append(data, model.UserLink{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
		})
	}

	return
}

func (m *Mem) CheckIsURLExists(longURL string) (string, error) {
	for long := range m.cacheMemory {
		if long == longURL {
			return m.cacheMemory[longURL], nil
		}
	}

	return "", nil
}

func (m *Mem) Close() error {
	return nil
}

func (m *Mem) Ping() bool {
	return m.cacheMemory == nil
}
