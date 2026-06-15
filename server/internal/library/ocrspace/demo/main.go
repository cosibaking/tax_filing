// OCR.space 模块独立测试服务
//
// 启动（在 server 目录下）：
//
//	go run ./internal/library/ocrspace/demo -config manifest/config/ocrspace.yaml
//
// 浏览器访问：http://localhost:8099
package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"xygo/internal/library/ocrspace"
)

//go:embed page.html
var pageFS embed.FS

const maxUploadSize = 10 << 20 // 10 MB

func main() {
	addr := flag.String("addr", ":8099", "监听地址")
	configPath := flag.String("config", "", "ocrspace 配置文件路径")
	apiKey := flag.String("apikey", "", "API Key（覆盖配置文件，测试可用 helloworld）")
	flag.Parse()

	client, cfgPath, err := newClient(*configPath, *apiKey)
	if err != nil {
		log.Fatalf("init ocrspace client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", servePage)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":         true,
			"configPath": cfgPath,
		})
	})
	mux.HandleFunc("POST /api/ocr/file", func(w http.ResponseWriter, r *http.Request) {
		handleOCRFile(w, r, client)
	})
	mux.HandleFunc("POST /api/ocr/url", func(w http.ResponseWriter, r *http.Request) {
		handleOCRURL(w, r, client)
	})

	log.Printf("OCR demo server listening on http://localhost%s", strings.TrimPrefix(*addr, ":"))
	log.Printf("config: %s", cfgPath)
	if err = http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func newClient(configPath, apiKey string) (*ocrspace.Client, string, error) {
	var cfg ocrspace.Config
	var usedPath string
	var err error

	if configPath != "" {
		cfg, err = ocrspace.LoadConfig(configPath)
		usedPath = configPath
	} else if apiKey == "" {
		cfg, err = ocrspace.LoadConfig("")
		usedPath = "manifest/config/ocrspace.yaml"
	} else {
		cfg = ocrspace.DefaultConfig()
		usedPath = "(命令行 -apikey)"
	}
	if err != nil && apiKey == "" {
		return nil, usedPath, err
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	if err = cfg.Validate(); err != nil {
		return nil, usedPath, err
	}
	client, err := ocrspace.NewClient(cfg)
	return client, usedPath, err
}

func servePage(w http.ResponseWriter, _ *http.Request) {
	data, err := pageFS.ReadFile("page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

type requestOptions struct {
	Language            string `json:"language"`
	Engine              int    `json:"engine"`
	Scale               bool   `json:"scale"`
	IsTable             bool   `json:"isTable"`
	DetectOrientation   bool   `json:"detectOrientation"`
	IsOverlayRequired   bool   `json:"isOverlayRequired"`
}

type urlOCRRequest struct {
	URL     string         `json:"url"`
	Options requestOptions `json:"options"`
}

func parseRequestOptions(raw requestOptions) *ocrspace.ParseOptions {
	opts := &ocrspace.ParseOptions{
		Language:          ocrspace.Language(raw.Language),
		Scale:             raw.Scale,
		IsTable:           raw.IsTable,
		DetectOrientation: raw.DetectOrientation,
		IsOverlayRequired: raw.IsOverlayRequired,
	}
	if raw.Engine >= 1 && raw.Engine <= 3 {
		opts.OCREngine = ocrspace.OCREngine(raw.Engine)
	}
	return opts
}

func handleOCRFile(w http.ResponseWriter, r *http.Request, client *ocrspace.Client) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeJSON(w, http.StatusBadRequest, fail(err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fail(fmt.Errorf("缺少 file 字段")))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fail(err))
		return
	}

	var opts requestOptions
	if raw := r.FormValue("options"); raw != "" {
		_ = json.Unmarshal([]byte(raw), &opts)
	}

	ctx, cancel := context.WithTimeout(r.Context(), client.Config().Timeout)
	defer cancel()

	resp, err := client.ParseBytes(ctx, data, header.Filename, parseRequestOptions(opts))
	if err != nil {
		writeJSON(w, http.StatusBadGateway, fail(err))
		return
	}
	writeOCRResponse(w, resp)
}

func handleOCRURL(w http.ResponseWriter, r *http.Request, client *ocrspace.Client) {
	var req urlOCRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, fail(err))
		return
	}
	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		writeJSON(w, http.StatusBadRequest, fail(fmt.Errorf("url 不能为空")))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), client.Config().Timeout)
	defer cancel()

	resp, err := client.ParseURL(ctx, req.URL, parseRequestOptions(req.Options))
	if err != nil {
		writeJSON(w, http.StatusBadGateway, fail(err))
		return
	}
	writeOCRResponse(w, resp)
}

func writeOCRResponse(w http.ResponseWriter, resp *ocrspace.ParseResponse) {
	if err := resp.Err(); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
			"raw":     resp,
			"text":    resp.CombinedText(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"text":         resp.CombinedText(),
		"exitCode":     resp.OCRExitCode,
		"pageCount":    len(resp.ParsedResults),
		"processingMs": resp.ProcessingTimeInMilliseconds,
		"raw":          resp,
	})
}

func fail(err error) map[string]any {
	return map[string]any{
		"success": false,
		"error":   err.Error(),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}
