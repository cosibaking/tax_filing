// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package audit

import (
	"context"
	"os"
	"testing"

	"xygo/internal/library/compliancecrypto"
)

const testKey = "0123456789abcdef0123456789abcdef"

func TestEncryptDecrypt(t *testing.T) {
	ctx := context.Background()
	t.Setenv("COMPLIANCE_ENCRYPTION_KEY", testKey)

	plaintext := "110101199001011234"
	encrypted, err := compliancecrypto.Encrypt(ctx, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("Encrypt() should not return plaintext")
	}

	decrypted, err := compliancecrypto.Decrypt(ctx, encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestMaskIdCard(t *testing.T) {
	got := compliancecrypto.MaskIdCard("110101199001011234")
	want := "110***********1234"
	if got != want {
		t.Errorf("MaskIdCard() = %q, want %q", got, want)
	}
}

func TestMaskBankAccount(t *testing.T) {
	got := compliancecrypto.MaskBankAccount("6222021234567890123")
	want := "****0123"
	if got != want {
		t.Errorf("MaskBankAccount() = %q, want %q", got, want)
	}
}

func TestEncryptWithoutKey(t *testing.T) {
	ctx := context.Background()
	os.Unsetenv("COMPLIANCE_ENCRYPTION_KEY")

	_, err := compliancecrypto.Encrypt(ctx, "test")
	if err == nil {
		t.Fatal("Encrypt() should fail when key is not configured")
	}
}
