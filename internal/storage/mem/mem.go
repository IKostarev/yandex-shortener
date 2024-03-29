package mem

import (
	"github.com/IKostarev/yandex-go-dev/internal/utils"
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

func (m *Mem) Save(long, corrID string, _ string) (string, error) {
	short := utils.RandomString()

	m.cacheMemory[short] = long
	m.cacheCorrelation[corrID] = long

	return short, nil
}

func (m *Mem) Get(short, corrID string, _ string) (string, string) {
	return m.cacheMemory[short], corrID
}

func (m *Mem) CheckIsURLExists(longURL string) (string, error) {
	for long := range m.cacheMemory {
		if long == longURL {
			return m.cacheMemory[longURL], nil
		}
	}

	return "", nil
}

func (m *Mem) IsDel(_ string) bool {
	return false
}

func (m *Mem) GetAllURLs(_ string) ([]string, string) {
	return nil, ""
}

func (m *Mem) GetAllShortURLs(_ string) ([]string, string) {
	return nil, ""
}

func (m *Mem) DeleteURL(_ []string, _ string) bool {
	return true
}

func (m *Mem) Close() error {
	return nil
}

func (m *Mem) Ping() bool {
	return m.cacheMemory == nil
}
