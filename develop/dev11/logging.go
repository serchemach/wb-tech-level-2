package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()

}

func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}

func (rww ResponseWriterWrapper) String() string {
	var b strings.Builder

	b.WriteString("{ \"headers\": {")
	cnt := 0
	for k, v := range (*rww.w).Header() {
		b.WriteString(fmt.Sprintf("\"%s\": \"%v\"", k, v))
		if cnt != len((*rww.w).Header())-1 {
			b.WriteString(",")
		}
		cnt++
	}

	b.WriteString(fmt.Sprintf("}, \"status_code\": %d,", *(rww.statusCode)))

	b.WriteString("\"body\":\"")
	b.WriteString(rww.body.String())
	b.WriteString("\" }")
	return b.String()
}

type LoggerMiddleware struct {
	handler http.Handler
	logger  *log.Logger
}

func (m *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rww := NewResponseWriterWrapper(w)
	start := time.Now()
	m.handler.ServeHTTP(rww, r)

	m.logger.Printf("{\"request\": \"%s %s\", \"response\": %s, \"time_taken\": \"%s\"}\n", r.Method, r.URL.Path, rww.String(), time.Since(start).String())
}

func NewLoggerMiddleware(handlerToWrap http.Handler, logger *log.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{handlerToWrap, logger}
}
