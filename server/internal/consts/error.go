// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package consts

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

// concealErrorSlice 需要对外隐藏真实错误、统一成友好提示的错误关键字
var concealErrorSlice = []string{ErrorORM}

// ErrorMessage 用于统一对外错误描述（非 debug 环境可使用）
// - 如果是我们标记的内部错误类型，则返回统一的“操作失败，请稍后重试！”
// - 否则直接返回 err.Error()
func ErrorMessage(err error) (message string) {
	if err == nil {
		return "操作失败！"
	}
	message = err.Error()
	for _, e := range concealErrorSlice {
		if gstr.Contains(message, e) {
			return "操作失败，请稍后重试！"
		}
	}
	return
}

// ApiErrorMessage 返回适合展示给 API 调用方的错误文案。
// 数据库错误默认隐藏 SQL，仅 debug 模式下返回底层原因。
func ApiErrorMessage(ctx context.Context, err error) string {
	if err == nil {
		return "操作失败"
	}

	current := gerror.Current(err)
	if current == nil {
		return "操作失败"
	}

	code := gerror.Code(err)
	if code == gcode.CodeNil {
		code = gerror.Code(current)
	}

	// 已显式设置业务错误码时，直接使用业务消息。
	if code != gcode.CodeNil && code != gcode.CodeInternalError && code != gcode.CodeDbOperationError {
		return current.Error()
	}

	// GoFrame 数据库错误外层通常是完整 SQL，优先取底层 cause。
	if code == gcode.CodeDbOperationError || gstr.Contains(current.Error(), "UPDATE ") || gstr.Contains(current.Error(), "INSERT ") {
		if cause := gerror.Cause(err); cause != nil {
			causeMsg := cause.Error()
			if gstr.Contains(causeMsg, "Duplicate entry") || gstr.Contains(causeMsg, "duplicate key") {
				if gstr.Contains(causeMsg, "uk_mobile") || gstr.Contains(causeMsg, "mobile") {
					return "手机号已被其他账号使用"
				}
				if gstr.Contains(causeMsg, "uk_username") || gstr.Contains(causeMsg, "username") {
					return "用户名已被占用"
				}
			}
		}
		if g.Cfg().MustGet(ctx, "system.debug").Bool() {
			if cause := gerror.Cause(err); cause != nil {
				return cause.Error()
			}
		}
		return "保存失败，请稍后重试"
	}

	msg := current.Error()
	for _, e := range concealErrorSlice {
		if gstr.Contains(msg, e) {
			return "操作失败，请稍后重试"
		}
	}
	return msg
}
