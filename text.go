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
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// TextOptions 定义文本处理器的选项
type TextOptions struct {
	Text     string
	Position image.Point
	Font     font.Face
	Color    color.Color
	// 旋转角度（度数，顺时针方向）
	Angle float64
}

// DefaultTextOptions 默认文本选项
var DefaultTextOptions = TextOptions{
	Font:  basicfont.Face7x13,
	Color: color.Black,
}

// TextProcessor 实现文本处理器
type TextProcessor struct {
	Options TextOptions
}

// NewTextProcessor 创建新的文本处理器
func NewTextProcessor(opts TextOptions) *TextProcessor {
	if opts.Font == nil {
		opts.Font = DefaultTextOptions.Font
	}
	if opts.Color == nil {
		opts.Color = DefaultTextOptions.Color
	}
	return &TextProcessor{Options: opts}
}

// WithAngle 设置文本旋转角度
func (p *TextProcessor) WithAngle(angle float64) *TextProcessor {
	p.Options.Angle = angle
	return p
}

// Process 实现Processor接口
func (p *TextProcessor) Process(img image.Image) (image.Image, error) {
	// 获取原始图片尺寸
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建gg上下文
	dc := gg.NewContext(width, height)

	// 绘制原始图像
	dc.DrawImage(img, 0, 0)

	// 设置字体和颜色
	dc.SetFontFace(p.Options.Font)
	dc.SetColor(p.Options.Color)

	// 如果有旋转角度
	if p.Options.Angle != 0 {
		// 保存当前状态
		dc.Push()

		// 将角度转换为弧度
		angle := p.Options.Angle * math.Pi / 180.0

		// 移动到文本位置
		dc.Translate(float64(p.Options.Position.X), float64(p.Options.Position.Y))
		// 旋转指定角度
		dc.Rotate(angle)
		// 绘制文本（从原点开始）
		dc.DrawString(p.Options.Text, 0, 0)

		// 恢复状态
		dc.Pop()
	} else {
		// 无旋转，直接绘制文本
		dc.DrawString(p.Options.Text, float64(p.Options.Position.X), float64(p.Options.Position.Y))
	}

	return dc.Image(), nil
}
