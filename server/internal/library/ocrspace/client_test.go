package ocrspace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigFromBytes(t *testing.T) {
	yaml := []byte(`
ocrspace:
  apiKey: test-key
  timeout: 30s
  defaultLanguage: chs
  defaultEngine: 3
`)
	cfg, err := LoadConfigFromBytes(yaml)
	if err != nil {
		t.Fatalf("LoadConfigFromBytes: %v", err)
	}
	if cfg.APIKey != "test-key" {
		t.Fatalf("apiKey = %q, want test-key", cfg.APIKey)
	}
	if cfg.Timeout != 30*time.Second {
		t.Fatalf("timeout = %v, want 30s", cfg.Timeout)
	}
	if cfg.DefaultLanguage != LanguageChineseSimplified {
		t.Fatalf("defaultLanguage = %q", cfg.DefaultLanguage)
	}
	if cfg.DefaultEngine != Engine3 {
		t.Fatalf("defaultEngine = %d", cfg.DefaultEngine)
	}
}

func TestLoadConfigFromBytesMissingKey(t *testing.T) {
	yaml := []byte(`ocrspace: {}`)
	_, err := LoadConfigFromBytes(yaml)
	if err == nil {
		t.Fatal("expected error for missing apiKey")
	}
}

func TestParseResponseHelpers(t *testing.T) {
	resp := &ParseResponse{
		OCRExitCode: 1,
		ParsedResults: []ParsedResult{
			{ParsedText: "page one"},
			{ParsedText: "page two"},
		},
	}
	if !resp.IsSuccess() {
		t.Fatal("expected success")
	}
	text := resp.CombinedText()
	if !strings.Contains(text, "page one") || !strings.Contains(text, "page two") {
		t.Fatalf("CombinedText = %q", text)
	}
	if resp.Err() != nil {
		t.Fatalf("unexpected error: %v", resp.Err())
	}

	fail := &ParseResponse{OCRExitCode: 4, IsErroredOnProcessing: true}
	msg := "fatal error"
	fail.ErrorMessage = &msg
	if fail.IsSuccess() {
		t.Fatal("expected failure")
	}
	if fail.Err() == nil {
		t.Fatal("expected Err()")
	}
}

func TestClientParseURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("apikey") != "test-key" {
			t.Errorf("apikey header = %q", r.Header.Get("apikey"))
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		if r.FormValue("url") != "https://example.com/a.jpg" {
			t.Errorf("url = %q", r.FormValue("url"))
		}
		if r.FormValue("language") != "eng" {
			t.Errorf("language = %q", r.FormValue("language"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"ParsedResults":[{"ParsedText":"hello world","FileParseExitCode":1}],
			"OCRExitCode":1,
			"IsErroredOnProcessing":false
		}`))
	}))
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.APIKey = "test-key"
	cfg.Endpoint = srv.URL

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	resp, err := client.ParseURL(context.Background(), "https://example.com/a.jpg", nil)
	if err != nil {
		t.Fatalf("ParseURL: %v", err)
	}
	if resp.CombinedText() != "hello world" {
		t.Fatalf("text = %q", resp.CombinedText())
	}

	text, err := client.ExtractTextFromURL(context.Background(), "https://example.com/a.jpg", nil)
	if err != nil {
		t.Fatalf("ExtractTextFromURL: %v", err)
	}
	if text != "hello world" {
		t.Fatalf("text = %q", text)
	}
}

func TestClientParseFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile: %v", err)
		}
		defer file.Close()
		if header.Filename != "test.png" {
			t.Errorf("filename = %q", header.Filename)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"ParsedResults":[{"ParsedText":"file ocr","FileParseExitCode":1}],
			"OCRExitCode":1,
			"IsErroredOnProcessing":false
		}`))
	}))
	defer srv.Close()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.png")
	if err := os.WriteFile(path, []byte("fake-png"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := DefaultConfig()
	cfg.APIKey = "test-key"
	cfg.Endpoint = srv.URL

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	text, err := client.ExtractTextFromFile(context.Background(), path, nil)
	if err != nil {
		t.Fatalf("ExtractTextFromFile: %v", err)
	}
	if text != "file ocr" {
		t.Fatalf("text = %q", text)
	}
}

func TestClientParseBase64(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		b64 := r.FormValue("base64Image")
		if !strings.HasPrefix(b64, "data:image/jpeg;base64,") {
			t.Errorf("base64Image prefix unexpected: %q", truncate(b64, 40))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"ParsedResults":[{"ParsedText":"base64 ocr","FileParseExitCode":1}],
			"OCRExitCode":1,
			"IsErroredOnProcessing":false
		}`))
	}))
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.APIKey = "test-key"
	cfg.Endpoint = srv.URL

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	resp, err := client.ParseBase64(context.Background(), "data:image/jpeg;base64,abc123", nil)
	if err != nil {
		t.Fatalf("ParseBase64: %v", err)
	}
	if resp.CombinedText() != "base64 ocr" {
		t.Fatalf("text = %q", resp.CombinedText())
	}
}

func TestEnvAPIKeyOverride(t *testing.T) {
	t.Setenv(envAPIKey, "env-key")
	yaml := []byte(`ocrspace: { apiKey: file-key }`)
	cfg, err := LoadConfigFromBytes(yaml)
	if err != nil {
		t.Fatalf("LoadConfigFromBytes: %v", err)
	}
	if cfg.APIKey != "env-key" {
		t.Fatalf("apiKey = %q, want env-key", cfg.APIKey)
	}
}
