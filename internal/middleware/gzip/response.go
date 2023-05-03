package gzip

import (
	"compress/gzip"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Response(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //TODO http.StatusInternalServerError
			logger.Error("new writer level is error: ", err)
			return
		}
		defer func(gz *gzip.Writer) {
			err = gz.Close()
			if err != nil {
				logger.Error("gzip.Reponse gz.Close() failed: ", err)
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
