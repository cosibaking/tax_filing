// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package compliancecrypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	envEncryptionKey = "COMPLIANCE_ENCRYPTION_KEY"
	cfgEncryptionKey = "compliance.encryptionKey"
	aesKeySize       = 32
)

// Encrypt AES-256-GCM 加密，返回 iv:authTag:ciphertext（均为 base64）
func Encrypt(ctx context.Context, plaintext string) (string, error) {
	key, err := resolveKey(ctx)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", gerror.Wrap(err, "create aes cipher failed")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", gerror.Wrap(err, "create gcm failed")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", gerror.Wrap(err, "generate nonce failed")
	}

	sealed := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	tagSize := gcm.Overhead()
	ciphertext := sealed[:len(sealed)-tagSize]
	tag := sealed[len(sealed)-tagSize:]

	return strings.Join([]string{
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(tag),
		base64.StdEncoding.EncodeToString(ciphertext),
	}, ":"), nil
}

// Decrypt 解密 Encrypt 产出的密文
func Decrypt(ctx context.Context, encrypted string) (string, error) {
	key, err := resolveKey(ctx)
	if err != nil {
		return "", err
	}

	parts := strings.Split(encrypted, ":")
	if len(parts) != 3 {
		return "", gerror.New("invalid encrypted format, expected iv:authTag:ciphertext")
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", gerror.Wrap(err, "decode nonce failed")
	}
	tag, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", gerror.Wrap(err, "decode auth tag failed")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", gerror.Wrap(err, "decode ciphertext failed")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", gerror.Wrap(err, "create aes cipher failed")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", gerror.Wrap(err, "create gcm failed")
	}

	sealed := append(ciphertext, tag...)
	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", gerror.Wrap(err, "decrypt failed")
	}

	return string(plaintext), nil
}

// MaskIdCard 身份证脱敏: 110***********1234
func MaskIdCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:3] + strings.Repeat("*", len(idCard)-7) + idCard[len(idCard)-4:]
}

// MaskBankAccount 银行卡脱敏: ****1234
func MaskBankAccount(account string) string {
	if len(account) < 4 {
		return account
	}
	return "****" + account[len(account)-4:]
}

func resolveKey(ctx context.Context) ([]byte, error) {
	raw := os.Getenv(envEncryptionKey)
	if raw == "" {
		raw = g.Cfg().MustGet(ctx, cfgEncryptionKey, "").String()
	}
	if raw == "" {
		return nil, gerror.New("COMPLIANCE_ENCRYPTION_KEY not configured")
	}
	if len(raw) != aesKeySize {
		return nil, gerror.Newf("encryption key must be exactly %d bytes, got %d", aesKeySize, len(raw))
	}
	return []byte(raw), nil
}
