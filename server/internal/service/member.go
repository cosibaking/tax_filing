// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"xygo/internal/model"
	"xygo/internal/model/input/adminin"
	"xygo/internal/model/input/memberin"
)

type (
	IMemberAuth interface {
		// Login 会员登录
		Login(ctx context.Context, in *memberin.LoginInput) (out *memberin.LoginOutput, err error)
		// Register 会员注册
		Register(ctx context.Context, in *memberin.RegisterInput) (out *memberin.RegisterOutput, err error)
		// Logout 会员退出登录
		Logout(ctx context.Context, tokenStr string) (err error)
	}
	IMemberCheckin interface {
		// GetCheckinInfo 获取签到信息（7天日历+连续天数+今日状态）
		GetCheckinInfo(ctx context.Context, memberId uint64) (out *memberin.CheckinInfoOutput, err error)
		// DoCheckin 执行签到
		DoCheckin(ctx context.Context, memberId uint64) (out *memberin.DoCheckinOutput, err error)
	}
	IAdminMemberGroup interface {
		// List 获取会员分组列表
		List(ctx context.Context, in *adminin.MemberGroupListInp) (list []adminin.MemberGroupListItem, total int, err error)
		// Save 保存会员分组（新增/编辑）
		Save(ctx context.Context, in *adminin.MemberGroupSaveInp) (uint, error)
		// Delete 删除会员分组
		Delete(ctx context.Context, id uint64) error
	}
	IAdminMemberMenu interface {
		// Tree 获取会员菜单树
		Tree(ctx context.Context, in *adminin.MemberMenuTreeInp) ([]*adminin.MemberMenuTreeItem, error)
		// Save 保存会员菜单（新增/编辑）
		// 对齐 BuildAdmin 类型约束：
		// - menu_dir/menu/nav_user_menu → 强制 no_login_valid = 0
		// - route → 强制 menu_type = "tab"
		// - button → 清空 path/component/menu_type/url
		Save(ctx context.Context, in *adminin.MemberMenuSaveInp) (uint, error)
		// Delete 删除会员菜单
		Delete(ctx context.Context, id uint64) error
	}
	IFrontendNotice interface {
		// List 通知列表（带已读状态，只返回已发布的通知）
		List(ctx context.Context, memberId uint64, in *memberin.NoticeListInput) (out *memberin.NoticeListOutput, err error)
		// MarkRead 标记单条通知已读
		MarkRead(ctx context.Context, memberId uint64, noticeId uint64) error
		// MarkAllRead 全部通知标记已读
		MarkAllRead(ctx context.Context, memberId uint64) error
	}
	IMemberUser interface {
		// GetInfo 获取当前会员信息
		GetInfo(ctx context.Context, memberId uint64) (out *memberin.GetInfoOutput, err error)
		// UpdateProfile 更新会员资料
		UpdateProfile(ctx context.Context, memberId uint64, in *memberin.UpdateProfileInput) (err error)
		// ChangePassword 修改密码
		ChangePassword(ctx context.Context, memberId uint64, in *memberin.ChangePasswordInput) (err error)
		// VerifyPassword 验证会员登录密码
		VerifyPassword(ctx context.Context, memberId uint64, password string) (err error)
		// GetByUsername 根据用户名获取会员
		GetByUsername(ctx context.Context, username string) (out *model.MemberUser, err error)
		// GetMenusByGroupId 根据分组ID获取菜单列表（前台会员用）
		// 1. 查分组的 rules（逗号分隔的菜单ID，或 * 表示全部）
		// 2. 按 rules 过滤 member_menu 中 status=1 的菜单
		// 3. 如果未登录（groupId=0），只返回 no_login_valid=1 的菜单
		GetMenusByGroupId(ctx context.Context, groupId uint64) (menus []memberin.FrontendMenuItem, err error)
	}
)

var (
	localMemberAuth       IMemberAuth
	localMemberCheckin    IMemberCheckin
	localAdminMemberGroup IAdminMemberGroup
	localAdminMemberMenu  IAdminMemberMenu
	localFrontendNotice   IFrontendNotice
	localMemberUser       IMemberUser
)

func MemberAuth() IMemberAuth {
	if localMemberAuth == nil {
		panic("implement not found for interface IMemberAuth, forgot register?")
	}
	return localMemberAuth
}

func RegisterMemberAuth(i IMemberAuth) {
	localMemberAuth = i
}

func MemberCheckin() IMemberCheckin {
	if localMemberCheckin == nil {
		panic("implement not found for interface IMemberCheckin, forgot register?")
	}
	return localMemberCheckin
}

func RegisterMemberCheckin(i IMemberCheckin) {
	localMemberCheckin = i
}

func AdminMemberGroup() IAdminMemberGroup {
	if localAdminMemberGroup == nil {
		panic("implement not found for interface IAdminMemberGroup, forgot register?")
	}
	return localAdminMemberGroup
}

func RegisterAdminMemberGroup(i IAdminMemberGroup) {
	localAdminMemberGroup = i
}

func AdminMemberMenu() IAdminMemberMenu {
	if localAdminMemberMenu == nil {
		panic("implement not found for interface IAdminMemberMenu, forgot register?")
	}
	return localAdminMemberMenu
}

func RegisterAdminMemberMenu(i IAdminMemberMenu) {
	localAdminMemberMenu = i
}

func FrontendNotice() IFrontendNotice {
	if localFrontendNotice == nil {
		panic("implement not found for interface IFrontendNotice, forgot register?")
	}
	return localFrontendNotice
}

func RegisterFrontendNotice(i IFrontendNotice) {
	localFrontendNotice = i
}

func MemberUser() IMemberUser {
	if localMemberUser == nil {
		panic("implement not found for interface IMemberUser, forgot register?")
	}
	return localMemberUser
}

func RegisterMemberUser(i IMemberUser) {
	localMemberUser = i
}
