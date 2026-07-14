package complianceai

import (
	"strings"
	"testing"
)

func TestRedactSensitive(t *testing.T) {
	input := "张三 110101199001011234 手机13812345678 卡号6222021234567890123"
	got := RedactSensitive(input)
	for _, secret := range []string{"110101199001011234", "13812345678", "6222021234567890123"} {
		if strings.Contains(got, secret) {
			t.Fatalf("secret leaked: %s", secret)
		}
	}
}
