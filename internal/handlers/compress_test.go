package handlers

//
//import (
//	"bytes"
//	"github.com/IKostarev/yandex-go-dev/internal/config"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//type mockStorage struct {
//	saveReturn string
//	saveErr    error
//	storage    map[string]string
//}
//
//func (m *mockStorage) Get(s, _ string) (string, string) {
//	return m.storage[s], ""
//}
//
//func (m *mockStorage) Close() error {
//	return nil
//}
//
//func (m *mockStorage) CheckIsURLExists(string) (string, error) {
//	return "", nil
//}
//
//func (m *mockStorage) Save(_, _ string) (string, error) {
//	return m.saveReturn, m.saveErr
//}
//
//func (m *mockStorage) Ping() bool {
//	return false
//}
//
//func TestApp_CompressHandler(t *testing.T) {
//	tests := []struct {
//		name       string
//		inputBody  []byte
//		shortURL   string
//		longURL    string
//		statusCode int
//		saveErr    bool
//		want       []byte
//	}{
//		{
//			name:       "good test",
//			inputBody:  []byte("http://test.site.com"),
//			shortURL:   "asdfghjk",
//			longURL:    "http://localhost:8080/asdfghjk",
//			statusCode: 201,
//			want:       []byte("http://localhost:8080/asdfghjk"),
//		},
//		{
//			name:       "body nil test",
//			inputBody:  []byte(""),
//			shortURL:   "",
//			longURL:    "http://localhost:8080",
//			statusCode: 400,
//			want:       []byte(""),
//		},
//		{
//			name:       "finally url bad test",
//			inputBody:  []byte("http://test.site.com"),
//			shortURL:   "asdfghjk",
//			longURL:    "http://localhost:8080",
//			statusCode: 400,
//			want:       []byte("http://localhost:8080/vghbrtyu"),
//		},
//		{
//			name:       "storage save error test",
//			inputBody:  []byte("http://test.site.com"),
//			shortURL:   "",
//			longURL:    "http://localhost:8080",
//			saveErr:    true,
//			statusCode: 400,
//			want:       []byte(""),
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				a := &App{
//					Config: config.Config{
//						BaseShortURL: "http://localhost:8080",
//					},
//					Storage: &mockStorage{
//						saveReturn: test.shortURL,
//						saveErr:    nil,
//
//					},
//				}
//				a.CompressHandler(w, r)
//			}))
//
//			defer ts.Close()
//
//			req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(test.inputBody))
//			if err != nil {
//				t.Fatalf("Failed to create request: %v", err)
//			}
//
//			resp, err := ts.Client().Do(req)
//			if err != nil {
//				t.Fatalf("Failed to make request: %v", err)
//			}
//			defer resp.Body.Close()
//
//			got, err := io.ReadAll(resp.Body)
//			if err != nil {
//				t.Fatalf("Failed to read response body: %v", err)
//			}
//			if string(got) != string(test.want) {
//				if test.statusCode != 400 {
//					t.Fatalf("Failed to status code, want: %v", test.statusCode)
//				}
//			}
//		})
//	}
//}
