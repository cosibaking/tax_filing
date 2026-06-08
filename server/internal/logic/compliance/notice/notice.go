package notice

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/dao"
)

// SendToMember 向指定会员发送系统通知（通过 MemberNotice，按分组投递）
func SendToMember(ctx context.Context, memberId uint64, title, content string) error {
	if memberId == 0 {
		return gerror.New("会员ID无效")
	}
	var member struct {
		GroupId uint64 `json:"group_id"`
	}
	err := g.DB().Model("xy_member").Ctx(ctx).Where("id", memberId).Fields("group_id").Scan(&member)
	if err != nil {
		return gerror.Wrap(err, "查询会员失败")
	}

	target := "all"
	targetId := uint64(0)
	if member.GroupId > 0 {
		target = "group"
		targetId = member.GroupId
	}

	_, err = dao.MemberNotice.Ctx(ctx).Data(g.Map{
		"title":      title,
		"content":    content,
		"type":       "system",
		"target":     target,
		"target_id":  targetId,
		"sender":     "合规服务",
		"status":     1,
		"created_at": time.Now().Unix(),
	}).Insert()
	return err
}
