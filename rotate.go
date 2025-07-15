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
)

// RotateProcessor 图像旋转处理器
// 对图像进行旋转处理，保持原有清晰度
type RotateProcessor struct {
	// 旋转角度（度数，顺时针方向）
	Angle float64
	// 背景颜色（旋转后可能出现的空白区域填充颜色）
	Background color.Color
	// 是否保持原始尺寸
	KeepSize bool
}

// Process 实现Processor接口
func (p *RotateProcessor) Process(img image.Image) (image.Image, error) {
	// 获取原始图片尺寸
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 将角度转换为弧度
	angle := p.Angle * math.Pi / 180.0

	// 计算旋转后的图像尺寸
	var width, height int
	if p.KeepSize {
		// 保持原始尺寸
		width = origWidth
		height = origHeight
	} else {
		// 计算旋转后的尺寸
		absCos := math.Abs(math.Cos(angle))
		absSin := math.Abs(math.Sin(angle))
		width = int(math.Ceil(float64(origWidth)*absCos + float64(origHeight)*absSin))
		height = int(math.Ceil(float64(origWidth)*absSin + float64(origHeight)*absCos))
	}

	// 创建gg上下文
	dc := gg.NewContext(width, height)

	// 设置背景颜色
	background := p.Background
	if background == nil {
		// 默认使用透明背景
		background = color.Transparent
	}

	// 填充背景
	dc.SetColor(background)
	dc.Clear()

	// 计算旋转中心点
	origCenterX := float64(origWidth) / 2.0
	origCenterY := float64(origHeight) / 2.0
	newCenterX := float64(width) / 2.0
	newCenterY := float64(height) / 2.0

	// 保存当前状态
	dc.Push()

	// 移动到新图像中心
	dc.Translate(newCenterX, newCenterY)
	// 旋转指定角度
	dc.Rotate(angle)
	// 移动回原点，考虑原图像尺寸
	dc.Translate(-origCenterX, -origCenterY)

	// 绘制原始图像
	dc.DrawImage(img, 0, 0)

	// 恢复状态
	dc.Pop()

	// 返回结果图像
	return dc.Image(), nil
}

// gg库内部已经实现了高质量的图像旋转和插值算法，不再需要自定义实现

// NewRotateProcessor 创建新的旋转处理器
func NewRotateProcessor(angle float64) *RotateProcessor {
	return &RotateProcessor{
		Angle:      angle,
		Background: color.Transparent,
		KeepSize:   false,
	}
}

// WithBackground 设置背景颜色
func (p *RotateProcessor) WithBackground(background color.Color) *RotateProcessor {
	p.Background = background
	return p
}

// WithKeepSize 设置是否保持原始尺寸
func (p *RotateProcessor) WithKeepSize(keepSize bool) *RotateProcessor {
	p.KeepSize = keepSize
	return p
}
