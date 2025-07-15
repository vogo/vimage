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
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// WatermarkProcessor 水印处理器
type WatermarkProcessor struct {
	Text     string     // 水印文本
	FontSize float64    // 字体大小
	Color    color.RGBA // 水印颜色
	Opacity  float64    // 不透明度 (0-1)
	Position string     // 位置 ("center", "top-left", "bottom-right" 等)
	Rotation float64    // 旋转角度
	FontFace font.Face  // 字体
}

// Process 实现Processor接口
func (p *WatermarkProcessor) Process(img image.Image) (image.Image, error) {
	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建新的上下文
	dc := gg.NewContext(width, height)

	// 绘制原图
	dc.DrawImage(img, 0, 0)

	// 设置字体
	if p.FontFace != nil {
		dc.SetFontFace(p.FontFace)
	} else if defaultFont != nil {
		face := truetype.NewFace(defaultFont, &truetype.Options{Size: p.FontSize})
		dc.SetFontFace(face)
	} else {
		dc.SetFontFace(basicfont.Face7x13)
	}

	// 设置颜色和透明度
	dc.SetColor(color.RGBA{
		R: p.Color.R,
		G: p.Color.G,
		B: p.Color.B,
		A: uint8(float64(p.Color.A) * p.Opacity),
	})

	// 计算水印位置
	textWidth, textHeight := dc.MeasureString(p.Text)
	var x, y float64

	switch p.Position {
	case "top-left":
		x, y = 10, 10+textHeight
	case "top-right":
		x, y = float64(width)-textWidth-10, 10+textHeight
	case "bottom-left":
		x, y = 10, float64(height)-10
	case "bottom-right":
		x, y = float64(width)-textWidth-10, float64(height)-10
	default: // center
		x, y = float64(width)/2-textWidth/2, float64(height)/2+textHeight/2
	}

	// 应用旋转
	if p.Rotation != 0 {
		dc.RotateAbout(gg.Radians(p.Rotation), x+textWidth/2, y-textHeight/2)
	}

	// 绘制水印文本
	dc.DrawString(p.Text, x, y)

	return dc.Image(), nil
}

// NewWatermarkProcessor 创建新的水印处理器
func NewWatermarkProcessor(text string, fontSize float64, color color.RGBA, opacity float64, position string, rotation float64) *WatermarkProcessor {
	// 验证参数
	if opacity < 0 || opacity > 1 {
		opacity = 0.5 // 默认半透明
	}

	return &WatermarkProcessor{
		Text:     text,
		FontSize: fontSize,
		Color:    color,
		Opacity:  opacity,
		Position: position,
		Rotation: rotation,
	}
}
