package handlers

//
//import (
//	"context"
//	"fmt"
//	"github.com/go-chi/chi/v5"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestApp_GetURLHandler(t *testing.T) {
//	tests := []struct {
//		name       string
//		idURL      string
//		getURL     string
//		statusCode int
//		want       string
//	}{
//		{
//			name:       "good test",
//			idURL:      "qwertyui",
//			getURL:     "http://localhost:8080/qwertyui",
//			statusCode: 307,
//			want:       "http://test.site.com",
//		},
//		{
//			name:       "id url is empty test",
//			idURL:      "",
//			statusCode: 400,
//			want:       "",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := &App{
//				Storage: &mockStorage{
//					storage: map[string]string{
//						tt.idURL: tt.want,
//					},
//				},
//			}
//
//			w := httptest.NewRecorder()
//			r := httptest.NewRequest("GET", "/{key}", nil)
//
//			rctx := chi.NewRouteContext()
//			rctx.URLParams.Add("id", tt.idURL)
//
//			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
//
//			app.GetURLHandler(w, r)
//
//			if w.Code != tt.statusCode {
//				t.Errorf("handler returned wrong status code: got %v want %v",
//					w.Code, tt.statusCode)
//			}
//
//			fmt.Printf("%+v\n", w.Header().Get("Location"))
//
//			if tt.want != "" && w.Header().Get("Location") != tt.want {
//				t.Errorf("handler returned unexpected body: got %v want %v",
//					w.Body.String(), tt.want)
//			}
//		})
//	}
//}
