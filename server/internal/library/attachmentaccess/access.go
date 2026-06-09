package attachmentaccess

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"

	"xygo/internal/dao"
	"xygo/internal/model/entity"
)

const defaultSignedTTL = 30 * time.Minute

type attachmentRecord struct {
	Id       uint64 `json:"id"`
	Topic    string `json:"topic"`
	UserId   uint64 `json:"user_id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Mimetype string `json:"mimetype"`
}

func secret(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "auth.jwt.secret").String()
}

// Sign 生成附件访问签名
func Sign(ctx context.Context, memberId, fileId uint64, expires int64) string {
	payload := fmt.Sprintf("%d:%d:%d", memberId, fileId, expires)
	mac := hmac.New(sha256.New, []byte(secret(ctx)))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify 校验签名
func Verify(ctx context.Context, memberId, fileId uint64, expires int64, sign string) bool {
	if memberId == 0 || fileId == 0 || expires <= 0 || strings.TrimSpace(sign) == "" {
		return false
	}
	if time.Now().Unix() > expires {
		return false
	}
	expected := Sign(ctx, memberId, fileId, expires)
	return hmac.Equal([]byte(expected), []byte(strings.TrimSpace(sign)))
}

// BuildSignedURL 生成带签名的附件访问 URL（相对路径）
func BuildSignedURL(ctx context.Context, memberId, fileId uint64, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = defaultSignedTTL
	}
	att, err := loadAttachmentByID(ctx, fileId)
	if err != nil {
		return "", err
	}
	if att.UserId != memberId {
		return "", gerror.New("无权访问该附件")
	}
	expires := time.Now().Add(ttl).Unix()
	sign := Sign(ctx, memberId, fileId, expires)
	return fmt.Sprintf("%s?id=%d&expires=%d&sign=%s", att.Url, fileId, expires, sign), nil
}

// ServeProtected 防盗链静态附件处理器
func ServeProtected(r *ghttp.Request) {
	ctx := r.Context()
	urlPath := normalizeURLPath(r.URL.Path)
	if urlPath == "" {
		r.Response.WriteStatus(403)
		return
	}

	fileId := r.Get("id").Uint64()
	expires := r.Get("expires").Int64()
	sign := r.Get("sign").String()

	if fileId > 0 && sign != "" {
		att, err := loadAttachmentByID(ctx, fileId)
		if err == nil && att.Url == urlPath && Verify(ctx, att.UserId, fileId, expires, sign) {
			writeAttachmentFile(r, att)
			return
		}
	}

	att, err := loadAttachmentByURL(ctx, urlPath)
	if err != nil || att == nil {
		r.Response.WriteStatus(404)
		return
	}

	// 会员上传的合规资料必须携带有效签名
	if att.Topic == "member" {
		r.Response.WriteStatus(403)
		return
	}

	if !isSameOriginReferer(r) {
		r.Response.WriteStatus(403)
		return
	}

	writeAttachmentFile(r, att)
}

func writeAttachmentFile(r *ghttp.Request, att *attachmentRecord) {
	localPath := filepath.Join("resource", "public", strings.TrimPrefix(att.Url, "/"))
	if !gfile.Exists(localPath) {
		r.Response.WriteStatus(404)
		return
	}
	if att.Mimetype != "" {
		r.Response.Header().Set("Content-Type", att.Mimetype)
	}
	r.Response.Header().Set("Cache-Control", "private, max-age=300")
	r.Response.Header().Set("X-Content-Type-Options", "nosniff")
	r.Response.ServeFile(localPath)
}

func isSameOriginReferer(r *ghttp.Request) bool {
	referer := strings.TrimSpace(r.Header.Get("Referer"))
	if referer == "" {
		return false
	}
	ref, err := url.Parse(referer)
	if err != nil {
		return false
	}
	host := r.Host
	if ref.Host == host {
		return true
	}
	// 允许 localhost 开发环境端口差异
	refHost := strings.Split(ref.Host, ":")[0]
	reqHost := strings.Split(host, ":")[0]
	return refHost != "" && refHost == reqHost
}

func normalizeURLPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func loadAttachmentByID(ctx context.Context, fileId uint64) (*attachmentRecord, error) {
	var att attachmentRecord
	err := dao.SysAttachment.Ctx(ctx).Where("id", fileId).Scan(&att)
	if err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.New("附件不存在")
	}
	return &att, nil
}

func loadAttachmentByURL(ctx context.Context, urlPath string) (*attachmentRecord, error) {
	var att attachmentRecord
	err := dao.SysAttachment.Ctx(ctx).Where(dao.SysAttachment.Columns().Url, urlPath).Scan(&att)
	if err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.New("附件不存在")
	}
	return &att, nil
}

// LoadMeta 读取附件元信息
func LoadMeta(ctx context.Context, fileId uint64) (*entity.SysAttachment, error) {
	var att entity.SysAttachment
	err := dao.SysAttachment.Ctx(ctx).Where("id", fileId).Scan(&att)
	if err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.New("附件不存在")
	}
	return &att, nil
}

// ParseSignedQuery 解析签名参数
func ParseSignedQuery(idStr, expiresStr, sign string) (fileId uint64, expires int64, ok bool) {
	fileId, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || fileId == 0 {
		return 0, 0, false
	}
	expires, err = strconv.ParseInt(strings.TrimSpace(expiresStr), 10, 64)
	if err != nil || expires <= 0 {
		return 0, 0, false
	}
	if strings.TrimSpace(sign) == "" {
		return 0, 0, false
	}
	return fileId, expires, true
}
