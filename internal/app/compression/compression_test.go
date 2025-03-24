package compression

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_ResponseCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	server := httptest.NewServer(GzipMiddleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
	gzr, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err)
	decompressed, err := io.ReadAll(gzr)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(decompressed))
}

func TestGzipMiddleware_RequestDecompression(t *testing.T) {
	compressedData := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(compressedData)
	_, err := gzWriter.Write([]byte("test request"))
	assert.NoError(t, err)
	gzWriter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "test request", string(body))
		w.Write([]byte("OK"))
	})
	server := httptest.NewServer(GzipMiddleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL, bytes.NewReader(compressedData.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(respBody))
}

func TestGzipMiddleware_NoCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	server := httptest.NewServer(GzipMiddleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Empty(t, resp.Header.Get("Content-Encoding"))
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(body))
}
