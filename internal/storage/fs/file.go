package fs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/model"
	"github.com/IKostarev/yandex-go-dev/internal/utils"
	uuID "github.com/google/uuid"
	"os"
)

type Fs struct {
	fh               *os.File
	cacheURL         map[string]string
	cacheCorrelation map[string]string
	cacheByID        map[uuID.UUID]map[string]string
	count            int64
}

type URLData struct {
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

func NewFsFromFile(path string) (*Fs, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return NewFs(file)
}

func NewFs(file *os.File) (*Fs, error) {
	fs := &Fs{
		fh:               file,
		cacheURL:         make(map[string]string),
		cacheCorrelation: make(map[string]string),
		cacheByID:        make(map[uuID.UUID]map[string]string),
		count:            0,
	}

	urlData := &URLData{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		err := json.NewDecoder(bytes.NewReader([]byte(line))).Decode(&urlData)
		if err != nil {
			logger.Errorf("error json decode in NewFs: %s", err)
		}

		fs.cacheURL[urlData.ShortURL] = urlData.OriginalURL
		fs.cacheCorrelation[urlData.CorrelationID] = urlData.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner is error: %w", err)
	}

	return fs, nil
}

func (m *Fs) Save(long, corrID string) (string, error) {
	urlData := &URLData{
		UUID:          fmt.Sprintf("%d", m.count),
		ShortURL:      utils.RandomString(),
		CorrelationID: corrID,
		OriginalURL:   long,
	}

	jsonData, err := json.Marshal(urlData)
	if err != nil {
		return "", fmt.Errorf("cannot marshal json: %w", err)
	}

	_, err = m.fh.Write([]byte("\n"))
	if err != nil {
		return "", fmt.Errorf("cannot write to file: %w", err)
	}

	_, err = m.fh.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("cannot write to file: %w", err)
	}

	m.count++

	m.cacheURL[urlData.ShortURL] = urlData.OriginalURL
	m.cacheCorrelation[urlData.CorrelationID] = urlData.OriginalURL
	//m.cacheByID[user] = map[string]string{urlData.ShortURL: urlData.OriginalURL}

	return urlData.ShortURL, nil
}

func (m *Fs) Get(short, corrID string) (string, string) {
	//for id, urls := range m.cacheByID {
	//	if id == user {
	//		return urls[short], corrID
	//	}
	//}
	//
	//return "", ""
	return m.cacheURL[short], corrID
}

func (m *Fs) GetUserLinks(user uuID.UUID) (data []model.UserLink, err error) {
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

func (m *Fs) CheckIsURLExists(longURL string) (string, error) {
	for _, urls := range m.cacheByID {
		for short, long := range urls {
			if long == longURL {
				return short, nil
			}
		}
	}

	return "", nil
}

func (m *Fs) Close() error {
	return m.fh.Close()
}

func (m *Fs) Ping() bool {
	return m.fh == nil
}
