// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package member

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"

	"xygo/api/member"
	"xygo/internal/consts"
	"xygo/internal/dao"
	"xygo/internal/library/attachmentaccess"
	"xygo/internal/library/contexts"
	"xygo/internal/model/do"
	"xygo/internal/model/entity"
	"xygo/internal/model/input/memberin"
	"xygo/internal/service"
	"xygo/utility"
)

// UploadFile 会员端文件上传（图片/PDF，写入附件表供合规资料引用）
func (c *ControllerV1) UploadFile(ctx context.Context, req *member.UploadFileReq) (res *member.UploadFileRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	r := g.RequestFromCtx(ctx)
	upFile := r.GetUploadFile("file")
	if upFile == nil {
		return nil, gerror.New("未选择文件")
	}

	if upFile.Size > 10*1024*1024 {
		return nil, gerror.New("文件大小不能超过 10MB")
	}
	ext := strings.ToLower(filepath.Ext(upFile.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".pdf": true}
	if !allowed[ext] {
		return nil, gerror.New("仅支持 jpg/png/pdf 格式")
	}

	mimeMap := map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".gif": "image/gif", ".webp": "image/webp", ".pdf": "application/pdf",
	}
	mimetype := mimeMap[ext]
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	subdir := gtime.Now().Format("Ymd")
	name := uuid.New().String() + ext
	savePath := filepath.Join("resource", "public", "attachment", "upload", subdir, name)
	saveDir := filepath.Dir(savePath)
	if !gfile.Exists(saveDir) {
		_ = gfile.Mkdir(saveDir)
	}

	upFile.Filename = name
	_, err = upFile.Save(saveDir)
	if err != nil {
		return nil, gerror.Newf("保存文件失败: %v", err)
	}

	url := "/attachment/upload/" + subdir + "/" + name
	sha1sum := ""
	if f, openErr := upFile.Open(); openErr == nil {
		if b, readErr := io.ReadAll(f); readErr == nil {
			sum := sha1.Sum(b)
			sha1sum = hex.EncodeToString(sum[:])
		}
		_ = f.Close()
	}

	attachmentId, err := saveMemberAttachment(ctx, memberId, url, upFile.Filename, upFile.Size, mimetype, sha1sum)
	if err != nil {
		return nil, gerror.Wrap(err, "保存附件记录失败")
	}

	accessURL, err := attachmentaccess.BuildSignedURL(ctx, memberId, attachmentId, 0)
	if err != nil {
		return nil, gerror.Wrap(err, "生成附件访问链接失败")
	}

	return &member.UploadFileRes{
		Url:          accessURL,
		Name:         name,
		Size:         upFile.Size,
		AttachmentId: attachmentId,
	}, nil
}

// AttachmentAccessURL 刷新当前会员可访问的附件签名链接
func (c *ControllerV1) AttachmentAccessURL(ctx context.Context, req *member.AttachmentAccessURLReq) (res *member.AttachmentAccessURLRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}
	accessURL, err := attachmentaccess.BuildSignedURL(ctx, memberId, req.FileId, 0)
	if err != nil {
		return nil, err
	}
	return &member.AttachmentAccessURLRes{Url: accessURL}, nil
}

func saveMemberAttachment(ctx context.Context, memberId uint64, url, originalName string, size int64, mimetype, sha1sum string) (uint64, error) {
	now := uint(utility.NowUnix())
	_, err := dao.SysAttachment.Ctx(ctx).Data(do.SysAttachment{
		Topic:      "member",
		UserId:     memberId,
		Url:        url,
		Name:       originalName,
		Size:       uint64(size),
		Mimetype:   mimetype,
		Quote:      1,
		Storage:    "local",
		Sha1:       sha1sum,
		CreateTime: now,
		UpdateTime: now,
	}).Insert()
	if err != nil {
		return 0, err
	}
	var record entity.SysAttachment
	err = dao.SysAttachment.Ctx(ctx).
		Where(dao.SysAttachment.Columns().Url, url).
		Where(dao.SysAttachment.Columns().UserId, memberId).
		OrderDesc(dao.SysAttachment.Columns().Id).
		Scan(&record)
	if err != nil || record.Id == 0 {
		return 0, gerror.New("获取附件ID失败")
	}
	return record.Id, nil
}

// ==================== 签到 ====================

// GetCheckinInfo 获取签到信息
func (c *ControllerV1) GetCheckinInfo(ctx context.Context, req *member.CheckinInfoReq) (res *member.CheckinInfoRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	out, err := service.MemberCheckin().GetCheckinInfo(ctx, memberId)
	if err != nil {
		return nil, err
	}

	res = &member.CheckinInfoRes{
		ContinuousDays: out.ContinuousDays,
		TodayChecked:   out.TodayChecked,
		TodayScore:     out.TodayScore,
	}
	res.WeekDays = make([]member.CheckinDayItem, len(out.WeekDays))
	for i, d := range out.WeekDays {
		res.WeekDays[i] = member.CheckinDayItem{
			Date:    d.Date,
			Checked: d.Checked,
			Score:   d.Score,
		}
	}
	return res, nil
}

// DoCheckin 执行签到
func (c *ControllerV1) DoCheckin(ctx context.Context, req *member.DoCheckinReq) (res *member.DoCheckinRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	out, err := service.MemberCheckin().DoCheckin(ctx, memberId)
	if err != nil {
		return nil, err
	}

	return &member.DoCheckinRes{
		Score:          out.Score,
		ContinuousDays: out.ContinuousDays,
	}, nil
}

// ==================== 积分记录 ====================

// ScoreLogList 积分记录列表（直接调 DAO，避免与 admin CRUD 命名冲突）
func (c *ControllerV1) ScoreLogList(ctx context.Context, req *member.ScoreLogListReq) (res *member.ScoreLogListRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	model := dao.MemberScoreLog.Ctx(ctx).Where("member_id", memberId)
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	var list []member.ScoreLogItem
	err = model.Page(page, pageSize).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []member.ScoreLogItem{}
	}

	return &member.ScoreLogListRes{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    count,
	}, nil
}

// ==================== 余额记录 ====================

// MoneyLogList 余额记录列表（直接调 DAO，避免与 admin CRUD 命名冲突）
func (c *ControllerV1) MoneyLogList(ctx context.Context, req *member.MoneyLogListReq) (res *member.MoneyLogListRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	model := dao.MemberMoneyLog.Ctx(ctx).Where("member_id", memberId)
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	var list []member.MoneyLogItem
	err = model.Page(page, pageSize).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []member.MoneyLogItem{}
	}

	return &member.MoneyLogListRes{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    count,
	}, nil
}

// ==================== 系统通知 ====================

// NoticeList 通知列表
func (c *ControllerV1) NoticeList(ctx context.Context, req *member.NoticeListReq) (res *member.NoticeListRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	out, err := service.FrontendNotice().List(ctx, memberId, &memberin.NoticeListInput{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	res = &member.NoticeListRes{
		Page:     out.Page,
		PageSize: out.PageSize,
		Total:    out.Total,
		Unread:   out.Unread,
	}
	res.List = make([]member.NoticeItem, len(out.List))
	for i, item := range out.List {
		res.List[i] = member.NoticeItem{
			Id:        item.Id,
			Title:     item.Title,
			Content:   item.Content,
			Type:      item.Type,
			Sender:    item.Sender,
			IsRead:    item.IsRead,
			CreatedAt: item.CreatedAt,
		}
	}
	return res, nil
}

// NoticeRead 标记通知已读
func (c *ControllerV1) NoticeRead(ctx context.Context, req *member.NoticeReadReq) (res *member.NoticeReadRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	err = service.FrontendNotice().MarkRead(ctx, memberId, req.NoticeId)
	if err != nil {
		return nil, err
	}
	return &member.NoticeReadRes{}, nil
}

// NoticeReadAll 全部通知已读
func (c *ControllerV1) NoticeReadAll(ctx context.Context, req *member.NoticeReadAllReq) (res *member.NoticeReadAllRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	err = service.FrontendNotice().MarkAllRead(ctx, memberId)
	if err != nil {
		return nil, err
	}
	return &member.NoticeReadAllRes{}, nil
}

// GetInfo 获取当前会员信息
func (c *ControllerV1) GetInfo(ctx context.Context, req *member.GetInfoReq) (res *member.GetInfoRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	output, err := service.MemberUser().GetInfo(ctx, memberId)
	if err != nil {
		return nil, err
	}

	return &member.GetInfoRes{
		Id:       output.Id,
		Username: output.Username,
		Nickname: output.Nickname,
		Avatar:   output.Avatar,
		Mobile:   output.Mobile,
		Email:    output.Email,
		Gender:   output.Gender,
		Level:    output.Level,
		GroupId:  output.GroupId,
		Score:       output.Score,
		Money:       output.Money,
		LastLoginAt: output.LastLoginAt,
		LastLoginIp: output.LastLoginIp,
	}, nil
}

// UpdateProfile 更新会员资料
func (c *ControllerV1) UpdateProfile(ctx context.Context, req *member.UpdateProfileReq) (res *member.UpdateProfileRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	input := &memberin.UpdateProfileInput{
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
		Gender:   req.Gender,
		Birthday: req.Birthday,
		Email:    req.Email,
	}

	err = service.MemberUser().UpdateProfile(ctx, memberId, input)
	if err != nil {
		return nil, err
	}

	return &member.UpdateProfileRes{}, nil
}

// ChangePassword 修改密码
func (c *ControllerV1) ChangePassword(ctx context.Context, req *member.ChangePasswordReq) (res *member.ChangePasswordRes, err error) {
	memberId := contexts.GetMemberId(ctx)
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")
	}

	input := &memberin.ChangePasswordInput{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	err = service.MemberUser().ChangePassword(ctx, memberId, input)
	if err != nil {
		return nil, err
	}

	return &member.ChangePasswordRes{}, nil
}

// GetMenus 获取当前会员可用菜单（按分组权限过滤）
func (c *ControllerV1) GetMenus(ctx context.Context, req *member.GetMenusReq) (res *member.GetMenusRes, err error) {
	// 获取当前会员的分组ID
	groupId := contexts.GetMemberGroupId(ctx)

	// 获取菜单列表
	allMenus, err := service.MemberUser().GetMenusByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	res = &member.GetMenusRes{
		Menus: make([]member.MemberMenuItem, 0),
		Nav:   make([]member.MemberMenuItem, 0),
		Rules: make([]member.MemberMenuItem, 0),
	}

	// 按类型分类
	for _, m := range allMenus {
		item := member.MemberMenuItem{
			Id:              m.Id,
			Pid:             m.Pid,
			Title:           m.Title,
			Name:            m.Name,
			Path:            m.Path,
			Component:       m.Component,
			Icon:            m.Icon,
			MenuType:        m.MenuType,
			Url:             m.Url,
			Type:            m.Type,
			NavShowChildren: m.NavShowChildren,
			NoLoginValid:    m.NoLoginValid,
			Sort:            m.Sort,
		}
		switch m.Type {
		case "menu_dir", "menu":
			res.Menus = append(res.Menus, item)
		case "nav", "nav_user_menu":
			res.Nav = append(res.Nav, item)
		case "route":
			res.Rules = append(res.Rules, item)
		case "button":
			res.Rules = append(res.Rules, item)
		}
	}

	attachNavDropdownChildren(res.Nav, allMenus)

	return res, nil
}

func attachNavDropdownChildren(nav []member.MemberMenuItem, flat []memberin.FrontendMenuItem) {
	byPid := make(map[uint64][]memberin.FrontendMenuItem)
	for _, m := range flat {
		if m.Type == "menu" {
			byPid[m.Pid] = append(byPid[m.Pid], m)
		}
	}
	for i := range nav {
		if nav[i].Type != "nav" || nav[i].NavShowChildren != 1 {
			continue
		}
		kids := byPid[nav[i].Id]
		if len(kids) == 0 {
			continue
		}
		sort.Slice(kids, func(a, b int) bool {
			if kids[a].Sort != kids[b].Sort {
				return kids[a].Sort < kids[b].Sort
			}
			return kids[a].Id < kids[b].Id
		})
		ch := make([]member.MemberMenuItem, 0, len(kids))
		for _, k := range kids {
			ch = append(ch, frontendMenuItemFromFlat(k))
		}
		nav[i].Children = ch
	}
}

func frontendMenuItemFromFlat(m memberin.FrontendMenuItem) member.MemberMenuItem {
	return member.MemberMenuItem{
		Id:              m.Id,
		Pid:             m.Pid,
		Title:           m.Title,
		Name:            m.Name,
		Path:            m.Path,
		Component:       m.Component,
		Icon:            m.Icon,
		MenuType:        m.MenuType,
		Url:             m.Url,
		Type:            m.Type,
		NavShowChildren: m.NavShowChildren,
		NoLoginValid:    m.NoLoginValid,
		Sort:            m.Sort,
	}
}
