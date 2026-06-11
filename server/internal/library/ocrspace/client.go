package ocrspace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Client OCR.space API 客户端
type Client struct {
	cfg    Config
	http   *http.Client
}

// NewClient 使用给定配置创建客户端
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		cfg: cfg,
		http: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// NewClientFromFile 从配置文件创建客户端
func NewClientFromFile(path string) (*Client, error) {
	cfg, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	return NewClient(cfg)
}

// Config 返回当前配置副本
func (c *Client) Config() Config {
	return c.cfg
}

func (c *Client) applyDefaults(opts *ParseOptions) ParseOptions {
	if opts == nil {
		return ParseOptions{
			Language:  c.cfg.DefaultLanguage,
			OCREngine: c.cfg.DefaultEngine,
		}
	}
	out := *opts
	if out.Language == "" {
		out.Language = c.cfg.DefaultLanguage
	}
	if out.OCREngine == 0 {
		out.OCREngine = c.cfg.DefaultEngine
	}
	return out
}

func (c *Client) writeOptions(w *multipart.Writer, opts ParseOptions) error {
	fields := map[string]string{
		"language": string(opts.Language),
	}
	if opts.IsOverlayRequired {
		fields["isOverlayRequired"] = "true"
	}
	if opts.Filetype != "" {
		fields["filetype"] = string(opts.Filetype)
	}
	if opts.DetectOrientation {
		fields["detectOrientation"] = "true"
	}
	if opts.IsCreateSearchablePdf {
		fields["isCreateSearchablePdf"] = "true"
	}
	if opts.IsSearchablePdfHideTextLayer {
		fields["isSearchablePdfHideTextLayer"] = "true"
	}
	if opts.Scale {
		fields["scale"] = "true"
	}
	if opts.IsTable {
		fields["isTable"] = "true"
	}
	if opts.OCREngine > 0 {
		fields["OCREngine"] = strconv.Itoa(int(opts.OCREngine))
	}

	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) doPost(ctx context.Context, body io.Reader, contentType string) (*ParseResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.Endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: create request: %w", err)
	}
	req.Header.Set("apikey", c.cfg.APIKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ocrspace: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 512))
	}

	var result ParseResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ocrspace: decode response: %w", err)
	}
	return &result, nil
}

// ParseURL 通过远程 URL 识别图片或 PDF（POST）
func (c *Client) ParseURL(ctx context.Context, imageURL string, opts *ParseOptions) (*ParseResponse, error) {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return nil, fmt.Errorf("ocrspace: url is required")
	}

	o := c.applyDefaults(opts)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("url", imageURL); err != nil {
		return nil, err
	}
	if err := c.writeOptions(w, o); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return c.doPost(ctx, &buf, w.FormDataContentType())
}

// ParseURLGet 通过 GET 端点识别远程 URL（仅支持 URL 方式，参数较少）
func (c *Client) ParseURLGet(ctx context.Context, imageURL string, opts *ParseOptions) (*ParseResponse, error) {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return nil, fmt.Errorf("ocrspace: url is required")
	}

	o := c.applyDefaults(opts)
	q := url.Values{}
	q.Set("apikey", c.cfg.APIKey)
	q.Set("url", imageURL)
	if o.Language != "" {
		q.Set("language", string(o.Language))
	}
	if o.IsOverlayRequired {
		q.Set("isOverlayRequired", "true")
	}

	reqURL := c.cfg.EndpointGet + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: create GET request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: GET request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: read GET response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ocrspace: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 512))
	}

	var result ParseResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ocrspace: decode GET response: %w", err)
	}
	return &result, nil
}

// ParseFile 上传本地文件识别
func (c *Client) ParseFile(ctx context.Context, filePath string, opts *ParseOptions) (*ParseResponse, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, fmt.Errorf("ocrspace: file path is required")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ocrspace: open file: %w", err)
	}
	defer f.Close()

	return c.ParseReader(ctx, f, filepath.Base(filePath), opts)
}

// ParseBytes 上传字节内容识别
func (c *Client) ParseBytes(ctx context.Context, data []byte, filename string, opts *ParseOptions) (*ParseResponse, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("ocrspace: file data is empty")
	}
	return c.ParseReader(ctx, bytes.NewReader(data), filename, opts)
}

// ParseReader 通过 io.Reader 上传文件识别
func (c *Client) ParseReader(ctx context.Context, r io.Reader, filename string, opts *ParseOptions) (*ParseResponse, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "upload.bin"
	}

	o := c.applyDefaults(opts)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, r); err != nil {
		return nil, fmt.Errorf("ocrspace: write file part: %w", err)
	}
	if err = c.writeOptions(w, o); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return c.doPost(ctx, &buf, w.FormDataContentType())
}

// ParseBase64 通过 Base64 字符串识别（需含 data URI 前缀，如 data:image/jpeg;base64,...）
func (c *Client) ParseBase64(ctx context.Context, base64Image string, opts *ParseOptions) (*ParseResponse, error) {
	base64Image = strings.TrimSpace(base64Image)
	if base64Image == "" {
		return nil, fmt.Errorf("ocrspace: base64Image is required")
	}

	o := c.applyDefaults(opts)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("base64Image", base64Image); err != nil {
		return nil, err
	}
	if err := c.writeOptions(w, o); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return c.doPost(ctx, &buf, w.FormDataContentType())
}

// ExtractText 识别并返回合并文本，失败时返回 error
func (c *Client) ExtractTextFromURL(ctx context.Context, imageURL string, opts *ParseOptions) (string, error) {
	resp, err := c.ParseURL(ctx, imageURL, opts)
	if err != nil {
		return "", err
	}
	if err = resp.Err(); err != nil {
		return "", err
	}
	return resp.CombinedText(), nil
}

// ExtractTextFromFile 从本地文件识别并返回合并文本
func (c *Client) ExtractTextFromFile(ctx context.Context, filePath string, opts *ParseOptions) (string, error) {
	resp, err := c.ParseFile(ctx, filePath, opts)
	if err != nil {
		return "", err
	}
	if err = resp.Err(); err != nil {
		return "", err
	}
	return resp.CombinedText(), nil
}

// ExtractTextFromBytes 从字节内容识别并返回合并文本
func (c *Client) ExtractTextFromBytes(ctx context.Context, data []byte, filename string, opts *ParseOptions) (string, error) {
	resp, err := c.ParseBytes(ctx, data, filename, opts)
	if err != nil {
		return "", err
	}
	if err = resp.Err(); err != nil {
		return "", err
	}
	return resp.CombinedText(), nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
