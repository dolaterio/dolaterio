package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func loggingHandler(h http.Handler) http.Handler {
	if log.Level < logrus.InfoLevel {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.RequestURI()
		log.WithFields(logrus.Fields{
			"method":        r.Method,
			"contentLength": r.ContentLength,
			"uri":           uri,
			"host":          r.Host,
		}).Info("REQ: ", r.Method, " ", uri)
		loggerWrapper := &responseLogger{w: w}
		h.ServeHTTP(loggerWrapper, r)
		log.WithFields(logrus.Fields{
			"status":        loggerWrapper.status,
			"contentLength": loggerWrapper.size,
		}).Info("RES: ", r.Method, " ", r.URL.Path)
	})
}
