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
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

// concealErrorSlice 需要对外隐藏真实错误、统一成友好提示的错误关键字
var concealErrorSlice = []string{ErrorORM}

// dbTechnicalPatterns 数据库底层错误特征（小写匹配）
var dbTechnicalPatterns = []string{
	"sql:",
	"no rows in result set",
	"select ",
	"insert ",
	"update ",
	"delete ",
	"duplicate entry",
	"duplicate key",
	"connection refused",
	"dial tcp",
	"pq:",
	"mysql:",
	"数据库执行异常",
}

// dbTechnicalSuffixes 包装错误中需剥离的技术后缀分隔符
var dbTechnicalSuffixes = []string{": sql:", ": SELECT ", ": INSERT ", ": UPDATE ", ": DELETE ", ": pq:", ": mysql:"}

// ErrorMessage 用于统一对外错误描述（非 debug 环境可使用）
// - 如果是我们标记的内部错误类型，则返回统一的“操作失败，请稍后重试！”
// - 否则直接返回 err.Error()
func ErrorMessage(err error) (message string) {
	if err == nil {
		return "操作失败！"
	}
	message = err.Error()
	if sanitized := sanitizeDbErrorMessage(message); sanitized != message {
		return sanitized
	}
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

	msg := current.Error()

	// 已显式设置业务错误码时，优先使用业务消息；若仍含数据库技术信息则脱敏。
	if code != gcode.CodeNil && code != gcode.CodeInternalError && code != gcode.CodeDbOperationError {
		if sanitized := sanitizeDbErrorMessage(msg); sanitized != msg {
			return sanitized
		}
		return msg
	}

	// GoFrame 数据库错误外层通常是完整 SQL，优先取底层 cause。
	if code == gcode.CodeDbOperationError || isDbTechnicalError(msg) {
		if cause := gerror.Cause(err); cause != nil {
			causeMsg := cause.Error()
			if friendly := friendlyDbMessage(causeMsg); friendly != "" {
				return friendly
			}
		}
		if g.Cfg().MustGet(ctx, "system.debug").Bool() {
			if cause := gerror.Cause(err); cause != nil {
				return cause.Error()
			}
		}
		if friendly := friendlyDbMessage(msg); friendly != "" {
			return friendly
		}
		return "操作失败，请稍后重试"
	}

	if sanitized := sanitizeDbErrorMessage(msg); sanitized != msg {
		return sanitized
	}

	for _, e := range concealErrorSlice {
		if gstr.Contains(msg, e) {
			return "操作失败，请稍后重试"
		}
	}
	return msg
}

// sanitizeDbErrorMessage 将含数据库技术细节的错误转为业务文案。
func sanitizeDbErrorMessage(msg string) string {
	if msg == "" {
		return "操作失败，请稍后重试"
	}
	if stripped := stripDbTechnicalSuffix(msg); stripped != msg {
		return stripped
	}
	if isDbTechnicalError(msg) {
		return friendlyDbMessage(msg)
	}
	return msg
}

func isDbTechnicalError(msg string) bool {
	lower := strings.ToLower(msg)
	for _, p := range dbTechnicalPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func stripDbTechnicalSuffix(msg string) string {
	for _, sep := range dbTechnicalSuffixes {
		if idx := gstr.Pos(msg, sep); idx > 0 {
			return gstr.Trim(msg[:idx])
		}
	}
	return msg
}

func friendlyDbMessage(msg string) string {
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "no rows in result set"):
		return ErrorNotData
	case strings.Contains(lower, "duplicate entry"), strings.Contains(lower, "duplicate key"):
		if strings.Contains(lower, "uk_mobile") || strings.Contains(lower, "mobile") {
			return "手机号已被其他账号使用"
		}
		if strings.Contains(lower, "uk_username") || strings.Contains(lower, "username") {
			return "用户名已被占用"
		}
		return "数据已存在"
	case strings.Contains(lower, "connection refused"), strings.Contains(lower, "dial tcp"):
		return "服务暂不可用，请稍后重试"
	case isDbTechnicalError(msg):
		return "操作失败，请稍后重试"
	default:
		return ""
	}
}
