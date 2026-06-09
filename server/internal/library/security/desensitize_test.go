package security

import "testing"

func TestMaskAddressUnicode(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"天堂之上", "天**上"},
		{"地上地下", "地**下"},
		{"北京市朝阳区建国路88号SOHO现代城A座1201", "北京*********************01"},
	}
	for _, c := range cases {
		got := MaskAddress(c.in)
		if got != c.want {
			t.Fatalf("MaskAddress(%q) = %q, want %q", c.in, got, c.want)
		}
		for _, r := range got {
			if r == '\uFFFD' {
				t.Fatalf("MaskAddress(%q) contains replacement char: %q", c.in, got)
			}
		}
	}
}
