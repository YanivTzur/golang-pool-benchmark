package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	testString = strings.Repeat("0", 1024)
	buffer = bytes.NewBufferString(testString)
)

func benchmarkHandler(b *testing.B, handler http.HandlerFunc) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("POST", "http://localhost:8080/", buffer)
			w := httptest.NewRecorder()
			handler(w, req)
		}
	})
}

func BenchmarkBasicHandler(b *testing.B) {
	benchmarkHandler(b, BasicHandler)
}

func BenchmarkObjectPoolBasedHandler(b *testing.B) {
	benchmarkHandler(b, ObjectPoolHandler)
}

func BenchmarkBufferPoolHandler(b *testing.B) {
	benchmarkHandler(b, BoundedPoolHandler)
}
