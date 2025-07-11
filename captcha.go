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
	Width      int        // 图片宽度
	Height     int        // 图片高度
	BgColor    color.RGBA // 背景颜色
	TextColor  color.RGBA // 文字颜色
	NoiseLines int        // 干扰线数量
	NoiseDots  int        // 干扰点数量
}

// DefaultCaptchaConfig 默认验证码配置
var DefaultCaptchaConfig = CaptchaConfig{
	Width:      120,
	Height:     40,
	BgColor:    color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色背景
	TextColor:  color.RGBA{R: 0, G: 0, B: 0, A: 255},       // 黑色文字
	NoiseLines: 5,
	NoiseDots:  50,
}

// GenCaptchaImage 生成验证码图片
// captcha: 验证码文本
// 返回: 图片字节缓冲区和错误信息
func GenCaptchaImage(captcha string) (*bytes.Buffer, error) {
	return GenCaptchaImageWithConfig(captcha, DefaultCaptchaConfig)
}

// GenCaptchaImageWithConfig 使用自定义配置生成验证码图片
func GenCaptchaImageWithConfig(captcha string, config CaptchaConfig) (*bytes.Buffer, error) {
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
func drawText(img *image.RGBA, text string, config CaptchaConfig) error {
	// 使用基础字体，通过调整间距实现更大的视觉效果
	face := basicfont.Face7x13

	// 计算文字总宽度和起始位置 (字符间距加大实现2倍大小的视觉效果)
	textWidth := len(text) * 16 // 字符间距加大到16像素 (原来7像素的2倍多)
	startX := (config.Width - textWidth) / 2
	startY := config.Height/2 + 8 // 垂直居中稍微偏下

	// 创建字体绘制器
	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{config.TextColor},
		Face: face,
	}

	// 逐个字符绘制，添加随机偏移
	for i, char := range text {
		// 添加随机垂直偏移
		yOffset := rand.Intn(8) - 4 // -4到4的随机偏移 (适应更大字体)
		// 添加随机水平间距
		xOffset := rand.Intn(6) - 3 // -3到3的随机偏移 (适应更大字体)

		x := startX + i*18 + xOffset // 字符间距约18像素 (原来的1.8倍)
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
func addNoiseLines(img *image.RGBA, config CaptchaConfig) {
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
func addNoiseDots(img *image.RGBA, config CaptchaConfig) {
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
