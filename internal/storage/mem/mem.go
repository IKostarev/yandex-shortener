package mem

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/model"
	"github.com/IKostarev/yandex-go-dev/internal/utils"
	uuID "github.com/google/uuid"
)

type Mem struct {
	cacheCorrelation map[string]string
	cacheByID        map[uuID.UUID]map[string]string
}

func NewMem() (*Mem, error) {
	m := &Mem{
		cacheCorrelation: make(map[string]string),
		cacheByID:        make(map[uuID.UUID]map[string]string),
	}

	return m, nil
}

func (m *Mem) Save(long, corrID string) (string, error) {
	short := utils.RandomString()

	fmt.Println("SAVE LONG = ", long)

	m.cacheCorrelation[corrID] = long
	//m.cacheByID[user] = map[string]string{short: long}

	return short, nil
}

func (m *Mem) Get(short, corrID string) (string, string) {

	for _, urls := range m.cacheByID {

		return urls[short], corrID

		//if id == user {
		//	fmt.Println("urls short = ", urls[short])
		//	return urls[short], corrID
		//}
	}

	return "", ""
}

func (m *Mem) GetUserLinks(user uuID.UUID) (data []model.UserLink, err error) {
	data = make([]model.UserLink, 0)

	for id, urls := range m.cacheByID {
		if id == user {
			for short, long := range urls {
				data = append(data, model.UserLink{
					OriginalURL: long,
					ShortURL:    short,
				})
			}
		}
	}

	return data, nil
}

func (m *Mem) CheckIsURLExists(longURL string) (string, error) {
	for _, urls := range m.cacheByID {
		for short, long := range urls {
			if long == longURL {
				return short, nil
			}
		}
	}

	return "", nil
}

func (m *Mem) Close() error {
	return nil
}

func (m *Mem) Ping() bool {
	return m.cacheByID == nil
}
