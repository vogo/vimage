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
	"unicode"

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
	// 最大文本宽度（像素）。>0 时按宽度自动换行
	MaxWidth float64
	// 行距倍数（相对于字体行高），用于换行模式
	LineSpacing float64
	// 文本对齐方式（左/中/右），用于换行模式
	Align gg.Align
	// 使用按字符换行（适合中文、日文等无空格语言）
	CharWrap bool
}

// DefaultTextOptions 默认文本选项
var DefaultTextOptions = TextOptions{
	Font:  basicfont.Face7x13,
	Color: color.Black,
	// 默认不限制宽度、不换行
	MaxWidth:    0,
	LineSpacing: 1.5,
	Align:       gg.AlignLeft,
	CharWrap:    false,
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
	if opts.LineSpacing == 0 {
		opts.LineSpacing = DefaultTextOptions.LineSpacing
	}
	if opts.Align == 0 { // gg.AlignLeft 的零值通常为 0，这里显式兜底
		opts.Align = DefaultTextOptions.Align
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
	ctx := NewImageProcessContext(img)

	err := p.ContextProcess(ctx)
	if err != nil {
		return nil, err
	}

	return ctx.dc.Image(), nil
}

// ContextProcess 实现 ContextProcessor 接口
func (p *TextProcessor) ContextProcess(ctx *ImageProcessContext) error {
	dc := ctx.DC()

	// 设置字体和颜色
	dc.SetFontFace(p.Options.Font)
	dc.SetColor(p.Options.Color)

	drawWrapped := p.Options.MaxWidth > 0

	// 如果需要宽度限制，并且文本包含 CJK（或显式启用 CharWrap），进行按字符换行预处理
	textToDraw := p.Options.Text
	if drawWrapped && (p.Options.CharWrap || containsCJK(textToDraw)) {
		textToDraw = wrapTextByRune(p.Options.Font, textToDraw, p.Options.MaxWidth)
	}

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
		if drawWrapped {
			// 使用锚点(ax, ay)为(0, 0)，表示从原点基线开始绘制
			dc.DrawStringWrapped(textToDraw, 0, 0, 0, 0, p.Options.MaxWidth, p.Options.LineSpacing, p.Options.Align)
		} else {
			dc.DrawString(textToDraw, 0, 0)
		}

		// 恢复状态
		dc.Pop()
	} else {
		// 无旋转
		if drawWrapped {
			dc.DrawStringWrapped(
				textToDraw,
				float64(p.Options.Position.X),
				float64(p.Options.Position.Y),
				0, // ax: 左对齐
				0, // ay: 基线对齐
				p.Options.MaxWidth,
				p.Options.LineSpacing,
				p.Options.Align,
			)
		} else {
			dc.DrawString(textToDraw, float64(p.Options.Position.X), float64(p.Options.Position.Y))
		}
	}

	return nil
}

// containsCJK 判断文本是否包含中日韩字符，用于决定是否采用按字符换行
func containsCJK(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) {
			return true
		}
	}
	return false
}

// wrapTextByRune 按字符宽度进行换行，适用于中文等无空格分词的文本
func wrapTextByRune(face font.Face, s string, maxWidth float64) string {
	// 支持原文本中的显式换行：分段分别处理
	segments := splitByNewline(s)
	lines := make([]string, 0, len(segments))

	d := font.Drawer{Face: face}

	for _, seg := range segments {
		current := ""
		for _, r := range seg {
			next := current + string(r)
			w := float64(d.MeasureString(next)) / 64.0
			if w <= maxWidth || current == "" {
				current = next
			} else {
				lines = append(lines, current)
				current = string(r)
			}
		}
		if current != "" {
			lines = append(lines, current)
		}
	}

	return joinWithNewline(lines)
}

func splitByNewline(s string) []string {
	out := []string{}
	cur := ""
	for _, r := range s {
		if r == '\n' {
			out = append(out, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	out = append(out, cur)
	return out
}

func joinWithNewline(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	// 手动拼接，避免引入 strings 以减少依赖
	out := lines[0]
	for i := 1; i < len(lines); i++ {
		out += "\n" + lines[i]
	}
	return out
}
