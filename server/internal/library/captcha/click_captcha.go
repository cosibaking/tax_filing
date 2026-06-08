// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

// =================================================================================
// 点选验证码（Click Captcha）
// 移植自 BuildAdmin 的 ClickCaptcha.php
// 存储方式：数据表 xy_captcha（非缓存，不依赖 Redis）
// 特性：生成时自动清理过期记录、校验后自动删除
// =================================================================================

package captcha

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/guid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// resOpen 优先从内嵌资源读取，找不到则走文件系统
func resOpen(path string) (io.ReadCloser, error) {
	if file := gres.Get(path); file != nil {
		return io.NopCloser(bytes.NewReader(file.Content())), nil
	}
	return os.Open(path)
}

// resReadFile 优先从内嵌资源读取文件内容
func resReadFile(path string) ([]byte, error) {
	if file := gres.Get(path); file != nil {
		return file.Content(), nil
	}
	return os.ReadFile(path)
}

// ==================== 配置 ====================

type ClickConfig struct {
	Mode          []string // text / icon
	Length        int      // 需要点击的数量
	ConfuseLength int      // 干扰元素数量
	Alpha         float64  // 透明度 0~100
	Expire        int64    // 过期时间（秒）
}

var DefaultClickConfig = ClickConfig{
	Mode:          []string{"text", "icon"},
	Length:        2,
	ConfuseLength: 2,
	Alpha:         70, // 提高透明度(0~100)让文字/图标更清晰可见
	Expire:        600,
}

// ==================== 图标字典 ====================

var iconDict = map[string]string{
	"aeroplane": "飞机", "apple": "苹果", "banana": "香蕉", "bell": "铃铛",
	"bicycle": "自行车", "bird": "小鸟", "bomb": "炸弹", "butterfly": "蝴蝶",
	"candy": "糖果", "crab": "螃蟹", "cup": "杯子", "dolphin": "海豚",
	"fire": "火", "guitar": "吉他", "hexagon": "六角形", "pear": "梨",
	"rocket": "火箭", "sailboat": "帆船", "snowflake": "雪花", "wolf head": "狼头",
}

var zhSet = []rune("们以我到他会作时要动国产的是工就年阶义发成部民可出能方进在和有大这主中为来分生对于学地用同行面说种过命度革而多子后自社加小机也经力线本电高量长党得实家定深法表着水理化争现所起好十战无农使前等反体合斗路图把结第里正新开论之物从当两些还天资事队点育重其思与间内去因件利相由压员气业代全组数果期导平各基或月然如应形想制心样都向变关问比展那它最及外没看治提五解系林者米群头意只明四道马认次文通但条较克又公孔领军流入席位情运器并飞原油放立题质指建区验活众很教决特此常石强极已根共直团统式转别造切九你取西持总料连任志观调么山程百报更见必真保热委手改管处己将修支识象先老光专什六型具示复安带每东增则完风回南广劳轮科北打积车计给节做务被整联步类集号列温装即毫知轴研单色坚据速防史拉世设达尔场织历花受求传断况采精金界品判参层止边清至万确究书")

// ==================== 数据结构 ====================

type clickPoint struct {
	Text   string `json:"text"`
	Icon   bool   `json:"icon"`
	Name   string `json:"name,omitempty"`
	Size   int    `json:"size"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type clickCaptchaData struct {
	Text   []clickPoint `json:"text"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
}

// ClickCaptchaResult 返回给前端
type ClickCaptchaResult struct {
	Id     string   `json:"id"`
	Text   []string `json:"text"`
	Base64 string   `json:"base64"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
}

// ==================== 数据库操作 ====================

const captchaTable = "xy_captcha"
const clickTolerance = 6 // 点击容差（像素），提升点选命中率

func md5Str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// cleanExpired 清理过期验证码（每次生成时调用，对齐 BuildAdmin）
func cleanExpired(ctx context.Context) {
	_, err := g.DB().Ctx(ctx).Exec(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE expire_time < ?", captchaTable),
		time.Now().Unix(),
	)
	if err != nil {
		g.Log().Debugf(ctx, "清理过期验证码失败: %v", err)
	}
}

// saveCaptcha 存储验证码到数据库（兼容 MySQL/PG）
func saveCaptcha(ctx context.Context, id string, textHints []string, captchaJSON string) error {
	now := time.Now().Unix()
	key := md5Str(id)
	code := md5Str(strings.Join(textHints, ","))

	// 先删旧记录再插入（替代 MySQL 的 REPLACE INTO，兼容 PG）
	_, _ = g.DB().Ctx(ctx).Model(captchaTable).Where("key", key).Delete()
	_, err := g.DB().Ctx(ctx).Model(captchaTable).Data(g.Map{
		"key":         key,
		"code":        code,
		"captcha":     captchaJSON,
		"create_time": now,
		"expire_time": now + DefaultClickConfig.Expire,
	}).Insert()
	if err != nil {
		g.Log().Errorf(ctx, "存储验证码失败: %v", err)
		return fmt.Errorf("验证码服务暂不可用，请检查数据库连接")
	}
	return nil
}

// loadCaptcha 从数据库读取验证码（兼容 MySQL/PG）
func loadCaptcha(ctx context.Context, id string) (captchaJSON string, ok bool) {
	key := md5Str(id)
	row, err := g.DB().Ctx(ctx).Model(captchaTable).
		Fields("captcha", "expire_time").
		Where("key", key).
		One()
	if err != nil || row.IsEmpty() {
		return "", false
	}
	if time.Now().Unix() > row["expire_time"].Int64() {
		// 已过期，删除
		deleteCaptcha(ctx, id)
		return "", false
	}
	return row["captcha"].String(), true
}

// deleteCaptcha 删除验证码记录（兼容 MySQL/PG）
func deleteCaptcha(ctx context.Context, id string) {
	key := md5Str(id)
	_, _ = g.DB().Ctx(ctx).Model(captchaTable).Where("key", key).Delete()
}

// ==================== 资源路径 ====================

func getResourceBase() string {
	return "resource/captcha"
}

// ==================== 公开 API ====================

// GenerateClick 生成点选验证码
func GenerateClick(ctx context.Context) (*ClickCaptchaResult, error) {
	cfg := DefaultClickConfig
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 清理过期记录（对齐 BuildAdmin 构造函数逻辑）
	cleanExpired(ctx)

	resBase := getResourceBase()

	// 1. 随机背景图
	bgFiles := []string{"click/bgs/1.png", "click/bgs/2.png", "click/bgs/3.png"}
	bgPath := filepath.Join(resBase, bgFiles[r.Intn(len(bgFiles))])
	bgFile, err := resOpen(bgPath)
	if err != nil {
		return nil, fmt.Errorf("打开背景图失败: %w", err)
	}
	defer bgFile.Close()

	bgImg, err := png.Decode(bgFile)
	if err != nil {
		return nil, fmt.Errorf("解码背景图失败: %w", err)
	}
	imgBounds := bgImg.Bounds()
	imgW, imgH := imgBounds.Dx(), imgBounds.Dy()

	// 2. 加载字体
	var fontFace font.Face
	fontPath := filepath.Join(resBase, "fonts/SourceHanSansCN-Normal.ttf")
	fontData, err := resReadFile(fontPath)
	if err != nil {
		g.Log().Warningf(ctx, "加载字体失败: %v，仅使用图标模式", err)
		cfg.Mode = []string{"icon"}
	} else {
		ft, parseErr := opentype.Parse(fontData)
		if parseErr == nil {
			fontFace, _ = opentype.NewFace(ft, &opentype.FaceOptions{
				Size: 32, DPI: 72, Hinting: font.HintingFull,
			})
		}
		if fontFace == nil {
			cfg.Mode = []string{"icon"}
		}
	}

	// 3. 生成随机点位
	totalLen := cfg.Length + cfg.ConfuseLength
	points := randClickPoints(cfg, totalLen, r, fontFace)

	// 4. 随机布局（碰撞检测）
	for i := range points {
		x, y := randPosition(points, imgW, imgH, points[i].Width, points[i].Height, points[i].Icon, r)
		points[i].X = x
		points[i].Y = y
	}

	// 5. 绘制
	canvas := image.NewRGBA(imgBounds)
	draw.Draw(canvas, imgBounds, bgImg, image.Point{}, draw.Src)
	for _, p := range points {
		if p.Icon {
			drawIcon(canvas, p, resBase, cfg.Alpha)
		} else if fontFace != nil {
			drawText(canvas, p, fontFace, cfg.Alpha)
		}
	}

	// 6. 只取前 Length 个为答案
	answerPoints := points[:cfg.Length]
	textHints := make([]string, cfg.Length)
	for i, p := range answerPoints {
		textHints[i] = p.Text
	}

	// 7. 编码图片
	var buf bytes.Buffer
	_ = png.Encode(&buf, canvas)
	base64Img := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	// 8. 存入数据库
	captchaId := guid.S()
	captchaDataObj := clickCaptchaData{Text: answerPoints, Width: imgW, Height: imgH}
	dataBytes, _ := json.Marshal(captchaDataObj)
	if err := saveCaptcha(ctx, captchaId, textHints, string(dataBytes)); err != nil {
		return nil, err
	}

	return &ClickCaptchaResult{
		Id: captchaId, Text: textHints, Base64: base64Img,
		Width: imgW, Height: imgH,
	}, nil
}

// VerifyClick 校验点选验证码（校验后自动删除，一次性）
func VerifyClick(ctx context.Context, captchaId string, info string) bool {
	return VerifyClickErr(ctx, captchaId, info) == nil
}

// VerifyClickErr 校验点选验证码，返回具体错误原因
func VerifyClickErr(ctx context.Context, captchaId string, info string) error {
	if captchaId == "" || info == "" {
		return fmt.Errorf("验证码参数缺失")
	}

	captchaJSON, ok := loadCaptcha(ctx, captchaId)
	if !ok {
		return fmt.Errorf("验证码已失效，请刷新后重试")
	}

	var captchaData clickCaptchaData
	if err := json.Unmarshal([]byte(captchaJSON), &captchaData); err != nil {
		deleteCaptcha(ctx, captchaId)
		return fmt.Errorf("验证码数据异常，请刷新后重试")
	}

	// 解析 info: "x1,y1-x2,y2;width;height"
	parts := strings.Split(info, ";")
	if len(parts) != 3 {
		deleteCaptcha(ctx, captchaId)
		return fmt.Errorf("验证码格式错误，请刷新后重试")
	}
	displayW, _ := strconv.ParseFloat(parts[1], 64)
	displayH, _ := strconv.ParseFloat(parts[2], 64)
	if displayW == 0 || displayH == 0 {
		deleteCaptcha(ctx, captchaId)
		return fmt.Errorf("验证码格式错误，请刷新后重试")
	}

	xPro := displayW / float64(captchaData.Width)
	yPro := displayH / float64(captchaData.Height)

	xyPairs := strings.Split(parts[0], "-")
	if len(xyPairs) != len(captchaData.Text) {
		deleteCaptcha(ctx, captchaId)
		return fmt.Errorf("验证码点击数量不正确，请刷新后重试")
	}

	for k, pair := range xyPairs {
		coords := strings.Split(pair, ",")
		if len(coords) != 2 {
			deleteCaptcha(ctx, captchaId)
			return fmt.Errorf("验证码格式错误，请刷新后重试")
		}
		clickX, _ := strconv.ParseFloat(coords[0], 64)
		clickY, _ := strconv.ParseFloat(coords[1], 64)
		point := captchaData.Text[k]

		// X 轴
		realX := clickX / xPro
		if realX < float64(point.X-clickTolerance) || realX > float64(point.X+point.Width+clickTolerance) {
			deleteCaptcha(ctx, captchaId)
			return fmt.Errorf("验证码错误，请重试")
		}
		// Y 轴（图标和文字坐标系不同）
		realY := clickY / yPro
		var phStart, phEnd float64
		if point.Icon {
			phStart = float64(point.Y)
			phEnd = float64(point.Y + point.Height)
		} else {
			phStart = float64(point.Y - point.Height)
			phEnd = float64(point.Y)
		}
		if realY < phStart-float64(clickTolerance) || realY > phEnd+float64(clickTolerance) {
			deleteCaptcha(ctx, captchaId)
			return fmt.Errorf("验证码错误，请重试")
		}
	}

	// 校验成功后再删除（对齐 BuildAdmin unset 逻辑）
	deleteCaptcha(ctx, captchaId)
	return nil
}

// ==================== 内部方法（与之前相同） ====================

func randClickPoints(cfg ClickConfig, length int, r *rand.Rand, fontFace font.Face) []clickPoint {
	var points []clickPoint
	if containsStr(cfg.Mode, "text") && fontFace != nil {
		for i := 0; i < length; i++ {
			ch := string(zhSet[r.Intn(len(zhSet))])
			fontSize := r.Intn(10) + 26 // 26~35，比之前的15~30明显增大
			points = append(points, clickPoint{
				Text: ch, Icon: false, Size: fontSize,
				Width: fontSize + 8, Height: fontSize + 8,
			})
		}
	}
	if containsStr(cfg.Mode, "icon") {
		iconKeys := make([]string, 0, len(iconDict))
		for k := range iconDict {
			iconKeys = append(iconKeys, k)
		}
		r.Shuffle(len(iconKeys), func(i, j int) { iconKeys[i], iconKeys[j] = iconKeys[j], iconKeys[i] })
		count := length
		if count > len(iconKeys) {
			count = len(iconKeys)
		}
		for i := 0; i < count; i++ {
			name := iconKeys[i]
			iconPath := filepath.Join(getResourceBase(), "click/icons/"+name+".png")
			w, h := getImageSize(iconPath)
			if w == 0 {
				w, h = 50, 50
			}
			// 确保图标不会太小，至少 50x50
			if w < 50 {
				h = h * 50 / w
				w = 50
			}
			if h < 50 {
				w = w * 50 / h
				h = 50
			}
			points = append(points, clickPoint{
				Text: fmt.Sprintf("<%s>", iconDict[name]), Icon: true, Name: name,
				Width: w, Height: h,
			})
		}
	}
	r.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
	if len(points) > length {
		points = points[:length]
	}
	return points
}

func randPosition(points []clickPoint, imgW, imgH, fontW, fontH int, isIcon bool, r *rand.Rand) (int, int) {
	for attempt := 0; attempt < 200; attempt++ {
		x := r.Intn(maxInt(1, imgW-fontW))
		y := fontH + r.Intn(maxInt(1, imgH-fontH*2))
		if checkPosition(points, x, y, fontW, fontH, isIcon) {
			return x, y
		}
	}
	return r.Intn(maxInt(1, imgW-fontW)), fontH + r.Intn(maxInt(1, imgH-fontH*2))
}

func checkPosition(points []clickPoint, x, y, w, h int, isIcon bool) bool {
	for _, v := range points {
		if v.X == 0 && v.Y == 0 {
			continue
		}
		if !((x+w) < v.X || x > (v.X+v.Width)) {
			curTop, curBottom := y, y+h
			if !isIcon {
				curTop, curBottom = y-h, y
			}
			hisTop, hisBottom := v.Y, v.Y+v.Height
			if !v.Icon {
				hisTop, hisBottom = v.Y-v.Height, v.Y
			}
			if !(curBottom < hisTop || curTop > hisBottom) {
				return false
			}
		}
	}
	return true
}

func drawIcon(canvas *image.RGBA, p clickPoint, resBase string, alpha float64) {
	iconPath := filepath.Join(resBase, "click/icons/"+p.Name+".png")
	iconFile, err := resOpen(iconPath)
	if err != nil {
		return
	}
	defer iconFile.Close()
	iconImg, err := png.Decode(iconFile)
	if err != nil {
		return
	}

	// 缩放图标到目标尺寸
	scaledIcon := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	srcBounds := iconImg.Bounds()
	for dy := 0; dy < p.Height; dy++ {
		for dx := 0; dx < p.Width; dx++ {
			sx := srcBounds.Min.X + dx*srcBounds.Dx()/p.Width
			sy := srcBounds.Min.Y + dy*srcBounds.Dy()/p.Height
			scaledIcon.Set(dx, dy, iconImg.At(sx, sy))
		}
	}

	alphaMask := image.NewUniform(color.Alpha{A: uint8(alpha * 255 / 100)})
	draw.DrawMask(canvas, image.Rect(p.X, p.Y, p.X+p.Width, p.Y+p.Height),
		scaledIcon, image.Point{}, alphaMask, image.Point{}, draw.Over)
}

func drawText(canvas *image.RGBA, p clickPoint, face font.Face, alpha float64) {
	a := uint8(alpha * 255 / 100)

	// 先画深色描边（增强对比度，在任何背景上都可见）
	strokeCol := color.NRGBA{R: 30, G: 30, B: 60, A: a}
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			sd := &font.Drawer{
				Dst: canvas, Src: image.NewUniform(strokeCol), Face: face,
				Dot: fixed.P(p.X+dx, p.Y+dy),
			}
			sd.DrawString(p.Text)
		}
	}

	// 再画亮色主体文字
	mainCol := color.NRGBA{R: 255, G: 255, B: 245, A: a}
	d := &font.Drawer{
		Dst: canvas, Src: image.NewUniform(mainCol), Face: face,
		Dot: fixed.P(p.X, p.Y),
	}
	d.DrawString(p.Text)
}

func getImageSize(path string) (int, int) {
	f, err := resOpen(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()
	cfg, err := png.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
