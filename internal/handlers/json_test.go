package handlers

//
//import (
//	"bytes"
//	"github.com/IKostarev/yandex-go-dev/internal/config"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestApp_JSONHandler(t *testing.T) {
//	tests := []struct {
//		name        string
//		body        []byte
//		statusCode  int
//		expectedURL string
//	}{
//		{
//			name:        "good request",
//			body:        []byte(`{"url": "http://test.site.com"}`),
//			statusCode:  http.StatusCreated,
//			expectedURL: "/qwertyui",
//		},
//		{
//			name:        "bad request",
//			body:        []byte(`{"invalid_json":`),
//			statusCode:  http.StatusBadRequest,
//			expectedURL: "",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				a := &App{
//					Config: config.Config{
//						BaseShortURL: tt.expectedURL,
//					},
//					Storage: &mockStorage{},
//				}
//				a.JSONHandler(w, r)
//			}))
//			defer ts.Close()
//
//			resp, err := http.Post(ts.URL+"/api/shorten", "application/json", bytes.NewBuffer(tt.body))
//			if err != nil {
//				t.Errorf("failed to make POST request: %v", err)
//			}
//			if resp.StatusCode != tt.statusCode {
//				t.Errorf("expected status code %v, got %v", tt.statusCode, resp.StatusCode)
//			}
//
//			defer resp.Body.Close()
//		})
//	}
//}
