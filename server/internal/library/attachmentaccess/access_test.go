package attachmentaccess

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func TestSignAndVerify(t *testing.T) {
	ctx := context.Background()
	_ = g.Cfg().MustGet(ctx, "auth.jwt.secret").String()

	memberId := uint64(1001)
	fileId := uint64(2002)
	expires := time.Now().Add(5 * time.Minute).Unix()
	sign := Sign(ctx, memberId, fileId, expires)
	if !Verify(ctx, memberId, fileId, expires, sign) {
		t.Fatal("expected valid signature")
	}
	if Verify(ctx, memberId, fileId, expires, "bad-sign") {
		t.Fatal("expected invalid signature to fail")
	}
	if Verify(ctx, memberId, fileId, time.Now().Add(-time.Minute).Unix(), sign) {
		t.Fatal("expected expired signature to fail")
	}
}
