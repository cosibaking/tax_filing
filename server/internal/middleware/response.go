// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"xygo/internal/consts"
)

// JsonResponse 统一的 API 返回结构。
type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"traceId,omitempty"`
}

// ResponseHandler 统一响应中间件：
// - 正常返回：code=0,message="ok",data=业务返回值,traceId=当前请求ID
// - 发生错误：code=错误码,message=错误消息,traceId=当前请求ID
func ResponseHandler(r *ghttp.Request) {
	// 先执行后续处理
	r.Middleware.Next()

	// 若已有响应内容且没有错误，认为业务自行完整处理了响应，不再包一层 JSON。
	if r.GetError() == nil && r.Response.BufferLength() > 0 {
		return
	}

	traceID := gctx.CtxId(r.Context())
	err := r.GetError()

	// 统一处理 HTTP 状态码：
	// - 对“已匹配到的业务路由”统一使用 200
	// - 保留真实 404/405 给“路由不存在/方法不允许”
	status := r.Response.Status
	if status == 0 {
		status = http.StatusOK
	}
	if status != http.StatusNotFound && status != http.StatusMethodNotAllowed {
		r.Response.WriteStatus(http.StatusOK)
	}

	// 成功场景：无 error，将 Handler 返回值作为 data
	if err == nil {
		resp := JsonResponse{
			Code:    consts.CodeOK.Code(),
			Message: "ok",
			Data:    r.GetHandlerResponse(),
			TraceID: traceID,
		}
		r.Response.ClearBuffer()
		r.Response.WriteJson(resp)
		return
	}

	// 失败场景：提取错误码与消息
	code := gerror.Code(err).Code()
	if code == consts.CodeFailed.Code() {
		// 未显式设置业务错误码时，使用统一的通用失败码
		code = consts.CodeFailed.Code()
	}

	resp := JsonResponse{
		Code:    code,
		Message: consts.ApiErrorMessage(r.Context(), err),
		TraceID: traceID,
	}

	r.Response.ClearBuffer()
	r.Response.WriteJson(resp)
}
