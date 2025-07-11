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
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TextOptions 定义文本处理器的选项
type TextOptions struct {
	Text     string
	Position image.Point
	Font     font.Face
	Color    color.Color
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

// Process 实现ImageProcessor接口
func (p *TextProcessor) Process(img image.Image) (image.Image, error) {
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(p.Options.Color),
		Face: p.Options.Font,
		Dot:  fixed.P(p.Options.Position.X, p.Options.Position.Y),
	}
	drawer.DrawString(p.Options.Text)

	return dst, nil
}
