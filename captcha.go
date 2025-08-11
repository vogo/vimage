/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vimage

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Width            int        // 图片宽度
	Height           int        // 图片高度
	BgColor          color.RGBA // 背景颜色
	TextColor        color.RGBA // 文字颜色
	NoiseLines       int        // 干扰线数量
	NoiseDots        int        // 干扰点数量
	Face             font.Face  // 字体
	CharSpacing      int        // 字符间距
	CharWidth        int        // 字符宽度
	CharYOffsetRange int        // 垂直随机偏移范围
	CharXOffsetRange int        // 水平随机偏移范围
}

// DefaultCaptchaConfig 默认验证码配置
var DefaultCaptchaConfig = &CaptchaConfig{
	Width:            120,                                        // 图片宽度
	Height:           40,                                         // 图片高度
	BgColor:          color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色背景
	TextColor:        color.RGBA{R: 0, G: 0, B: 0, A: 255},       // 黑色文字
	NoiseLines:       5,                                          // 干扰线数量
	NoiseDots:        50,                                         // 干扰点数量
	Face:             basicfont.Face7x13,                         // 字体
	CharSpacing:      18,                                         // 字符间距
	CharWidth:        16,                                         // 字符宽度
	CharYOffsetRange: 8,                                          // 字符垂直随机偏移范围
	CharXOffsetRange: 6,                                          // 字符水平随机偏移范围
}

// GenCaptchaImage 生成验证码图片
// captcha: 验证码文本
// 返回: 图片字节缓冲区和错误信息
func GenCaptchaImage(captcha string) (*bytes.Buffer, error) {
	return GenCaptchaImageWithConfig(captcha, DefaultCaptchaConfig)
}

// GenCaptchaImageWithConfig 使用自定义配置生成验证码图片
func GenCaptchaImageWithConfig(captcha string, config *CaptchaConfig) (*bytes.Buffer, error) {
	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 填充背景色
	draw.Draw(img, img.Bounds(), &image.Uniform{config.BgColor}, image.Point{}, draw.Src)

	// 添加干扰线
	addNoiseLines(img, config)

	// 绘制验证码文字
	if err := drawText(img, captcha, config); err != nil {
		return nil, err
	}

	// 添加干扰点
	addNoiseDots(img, config)

	// 将图片编码为PNG格式
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf, nil
}

// drawText 绘制验证码文字
func drawText(img *image.RGBA, text string, config *CaptchaConfig) error {
	// 使用配置中的参数，如果为0则使用默认值
	charSpacing := config.CharSpacing
	charWidth := config.CharWidth
	yOffsetRange := config.CharYOffsetRange
	xOffsetRange := config.CharXOffsetRange

	// 如果配置参数为0，则根据字体类型设置默认值
	if charSpacing == 0 {
		charSpacing = 45
	}
	if charWidth == 0 {
		charWidth = 40
	}
	if yOffsetRange == 0 {
		yOffsetRange = 12
	}
	if xOffsetRange == 0 {
		xOffsetRange = 8
	}

	// 计算文字总宽度和起始位置
	textWidth := len(text) * charWidth
	startX := (config.Width - textWidth) / 2
	startY := config.Height/2 + config.Height/8 // 根据图片高度动态调整垂直位置

	// 创建字体绘制器
	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{config.TextColor},
		Face: config.Face,
	}

	// 逐个字符绘制，添加随机偏移
	for i, char := range text {
		// 添加随机垂直偏移
		yOffset := rand.Intn(yOffsetRange*2) - yOffsetRange
		// 添加随机水平间距
		xOffset := rand.Intn(xOffsetRange*2) - xOffsetRange

		x := startX + i*charSpacing + xOffset
		y := startY + yOffset

		d.Dot = fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}

		d.DrawString(string(char))
	}

	return nil
}

// addNoiseLines 添加干扰线
func addNoiseLines(img *image.RGBA, config *CaptchaConfig) {
	for i := 0; i < config.NoiseLines; i++ {
		// 随机起点和终点
		x1 := rand.Intn(config.Width)
		y1 := rand.Intn(config.Height)
		x2 := rand.Intn(config.Width)
		y2 := rand.Intn(config.Height)

		// 随机颜色（较浅的灰色）
		lineColor := color.RGBA{
			R: uint8(rand.Intn(100) + 100), // 100-200之间
			G: uint8(rand.Intn(100) + 100),
			B: uint8(rand.Intn(100) + 100),
			A: 255,
		}

		// 绘制线条
		DrawLine(img, x1, y1, x2, y2, lineColor)
	}
}

// addNoiseDots 添加干扰点
func addNoiseDots(img *image.RGBA, config *CaptchaConfig) {
	for i := 0; i < config.NoiseDots; i++ {
		x := rand.Intn(config.Width)
		y := rand.Intn(config.Height)

		// 随机颜色（较浅）
		dotColor := color.RGBA{
			R: uint8(rand.Intn(150) + 100), // 100-250之间
			G: uint8(rand.Intn(150) + 100),
			B: uint8(rand.Intn(150) + 100),
			A: 255,
		}

		img.Set(x, y, dotColor)
	}
}
