// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package cmd

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"xygo/internal/controller/admin"
	"xygo/internal/controller/hello"
	"xygo/internal/controller/member"
	"xygo/internal/controller/site"
	"xygo/internal/controller/system"
	"xygo/internal/controller/wm"
	"xygo/internal/addon"
	"xygo/internal/library/attachmentaccess"
	"xygo/internal/library/cache"
	"xygo/internal/library/monitor"
	"xygo/internal/library/queue"
	cronlogic "xygo/internal/logic/cron"
	"xygo/internal/middleware"
	"xygo/internal/websocket"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// ✨ 初始化缓存系统（对齐 HotGo）
			cache.Init(ctx)

			// ✨ 初始化性能监控
			monitor.InitPerformanceMonitor(ctx)

			// ✨ 启动 WebSocket Hub
			websocket.Start()

			// ✨ 启动定时任务调度
			cronlogic.StartAll(ctx)

			// ✨ 初始化消息队列 & 启动消费者
			queue.Init(ctx)
			queue.StartConsumers(ctx)

			s := g.Server()

			// =============== 前台模板页面路由（GoFrame 模板渲染，SEO 友好） ===============
			// 暂时禁用纯 HTML 模板页面，恢复 SPA 默认首页
			// s.BindHandler("GET:/", site.PageIndex)
			// site.RegisterNavRoutes(s)

			// 静态文件服务（packed 内可能含旧 dist，部署时优先读磁盘）
			s.SetServerRoot("resource/public/dist")
			s.BindHandler("GET:/attachment/*", attachmentaccess.ServeProtected)
			s.AddStaticPath("/m", "resource/public/mobile")
			s.SetIndexFolder(false)

			// 优先从磁盘提供 dist（避免 packed 内嵌旧版 index.html / assets）
			s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				if r.Method != http.MethodGet && r.Method != http.MethodHead {
					return
				}
				path := r.URL.Path
				switch {
				case path == "/" || path == "/index.html":
					if serveDistIndexFromDisk(r) {
						r.ExitAll()
					}
				case strings.HasPrefix(path, "/assets/"):
					rel := strings.TrimPrefix(path, "/")
					if serveDistFileFromDisk(r, rel) {
						r.ExitAll()
					}
				}
			})

			// SPA 回退：非 API 路径返回 index.html，供 Hash 路由接管
			s.BindStatusHandler(http.StatusNotFound, func(r *ghttp.Request) {
				if r.Method != http.MethodGet && r.Method != http.MethodHead {
					return
				}
				if isBackendOrApiPath(r.URL.Path) {
					return
				}
				if serveDistIndexFromDisk(r) {
					r.ExitAll()
				}
			})

			// =============== 前端对接路由（受 CORS 白名单保护） ===============
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.CORS,            // CORS 白名单（配置化）
					middleware.SlowApiMonitor,  // 慢接口监控
					middleware.ResponseHandler, // 统一响应包装
				)

				// 示例接口（无鉴权）
				group.Bind(
					hello.NewV1(),
				)

				// 系统基础接口（无鉴权）
				group.Bind(
					system.NewV1(),
					site.NewV1(),
				)

			// 后台管理接口（带鉴权）
			group.Group("/", func(ag *ghttp.RouterGroup) {
				ag.Middleware(middleware.AdminAuth)
				ag.Middleware(middleware.AdminPermission)
				ag.Middleware(middleware.DemoGuard)
				ag.Middleware(middleware.OperationLog)
					ag.Bind(
						admin.NewV1(),
					)
				})

			// 会员接口（前台用户，使用 Xy-User-Token）
			group.Group("/member", func(mg *ghttp.RouterGroup) {
				mg.Middleware(middleware.MemberAuth)
				mg.Middleware(middleware.DemoGuard)
				mg.Bind(
					member.NewV1(),
				)
			})
			// 微信小程序接口（复用 MemberAuth 体系，/wm/auth/login 已在白名单）
			group.Group("/wm", func(wg *ghttp.RouterGroup) {
				wg.Middleware(middleware.MemberAuth)
				wg.Bind(
					wm.NewV1(),
				)
			})
			})

			// =============== WebSocket 端点 ===============
			s.Group("/socket", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.CORS)
				group.Middleware(middleware.WsAuth)
				group.GET("/", websocket.WsHandler)
			})

			// =============== 开放 API 路由（不受 CORS 白名单限制） ===============
			// 用于第三方系统对接（支付回调、订单推送、开放平台等）
			// 安全性通过 API Key / 签名验证保证，而非 CORS
			// s.Group("/api", func(group *ghttp.RouterGroup) {
			// 	group.Middleware(
			// 		middleware.CORSOpen,        // 开放跨域（允许任何来源）
			// 		middleware.ResponseHandler, // 统一响应包装
			// 	)
			// 	// TODO: 添加 API Key / 签名验证中间件
			// 	// group.Middleware(middleware.ApiAuth)
			// 	// group.Bind(openapi.NewV1())
			// })

			// =============== 扩展模块路由（由 addon installer 自动管理） ===============
			addon.MountAll(s)

			s.Run()
			return nil
		},
	}
)

const distRoot = "resource/public/dist"

func serveDistFileFromDisk(r *ghttp.Request, rel string) bool {
	localPath := gfile.Join(distRoot, rel)
	if !gfile.Exists(localPath) || gfile.IsDir(localPath) {
		return false
	}
	r.Response.ServeFile(localPath)
	return true
}

func serveDistIndexFromDisk(r *ghttp.Request) bool {
	localPath := gfile.Join(distRoot, "index.html")
	if !gfile.Exists(localPath) {
		return false
	}
	r.Response.ClearBuffer()
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.WriteStatus(http.StatusOK, gfile.GetContents(localPath))
	return true
}

func isBackendOrApiPath(path string) bool {
	prefixes := []string{
		"/admin", "/member", "/wm", "/site", "/system", "/socket",
		"/swagger", "/api", "/attachment", "/assets", "/captcha", "/hello", "/m/",
	}
	for _, prefix := range prefixes {
		if path == prefix || strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
