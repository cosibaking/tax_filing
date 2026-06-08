// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package system

import (
	"context"

	api "xygo/api/system"
	"xygo/internal/library/captcha"
)

// ClickCaptcha 获取点选验证码
func (c *ControllerV1) ClickCaptcha(ctx context.Context, req *api.ClickCaptchaReq) (res *api.ClickCaptchaRes, err error) {
	result, err := captcha.GenerateClick(ctx)
	if err != nil {
		return nil, err
	}
	return &api.ClickCaptchaRes{
		Id:     result.Id,
		Text:   result.Text,
		Base64: result.Base64,
		Width:  result.Width,
		Height: result.Height,
	}, nil
}

// CheckClickCaptcha 校验点选验证码
func (c *ControllerV1) CheckClickCaptcha(ctx context.Context, req *api.CheckClickCaptchaReq) (res *api.CheckClickCaptchaRes, err error) {
	if err = captcha.VerifyClickErr(ctx, req.Id, req.Info); err != nil {
		return nil, err
	}
	return &api.CheckClickCaptchaRes{}, nil
}
