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
	"fmt"
	"image"
)

// CutPosition 定义切割位置
type CutPosition string

const (
	// CutPositionCenter 居中切割
	CutPositionCenter CutPosition = "center"
	// CutPositionTop 从顶部切割
	CutPositionTop CutPosition = "top"
	// CutPositionBottom 从底部切割
	CutPositionBottom CutPosition = "bottom"
	// CutPositionLeft 从左侧切割
	CutPositionLeft CutPosition = "left"
	// CutPositionRight 从右侧切割
	CutPositionRight CutPosition = "right"
)

// CutProcessor 图像切割处理器
// 从原始图像中切割出指定区域
type CutProcessor struct {
	// 目标宽度和高度
	Width  int
	Height int
	// 切割位置
	Position CutPosition
	// 自定义切割区域（如果指定了，则忽略Position）
	X int // 左上角X坐标
	Y int // 左上角Y坐标
	// 是否使用自定义区域
	UseCustomRegion bool
}

// Process 实现ImageProcessor接口
func (p *CutProcessor) Process(img image.Image) (image.Image, error) {
	// 获取原始图片尺寸
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 验证目标尺寸
	if p.Width <= 0 || p.Height <= 0 {
		return nil, fmt.Errorf("无效的切割尺寸: %dx%d", p.Width, p.Height)
	}

	// 检查目标尺寸是否超过原始尺寸
	if p.Width > origWidth || p.Height > origHeight {
		return nil, fmt.Errorf("切割尺寸(%dx%d)超过原始尺寸(%dx%d)",
			p.Width, p.Height, origWidth, origHeight)
	}

	// 计算切割的起始位置
	var x, y int

	if p.UseCustomRegion {
		// 使用自定义区域
		x = p.X
		y = p.Y

		// 验证自定义区域是否有效
		if x < 0 || y < 0 || x+p.Width > origWidth || y+p.Height > origHeight {
			return nil, fmt.Errorf("无效的切割区域: 起点(%d,%d), 尺寸(%dx%d), 原始尺寸(%dx%d)",
				x, y, p.Width, p.Height, origWidth, origHeight)
		}
	} else {
		// 根据位置计算起始点
		switch p.Position {
		case CutPositionTop:
			// 从顶部开始，水平居中
			x = (origWidth - p.Width) / 2
			y = 0
		case CutPositionBottom:
			// 从底部开始，水平居中
			x = (origWidth - p.Width) / 2
			y = origHeight - p.Height
		case CutPositionLeft:
			// 从左侧开始，垂直居中
			x = 0
			y = (origHeight - p.Height) / 2
		case CutPositionRight:
			// 从右侧开始，垂直居中
			x = origWidth - p.Width
			y = (origHeight - p.Height) / 2
		default: // CutPositionCenter 或其他值
			// 居中切割
			x = (origWidth - p.Width) / 2
			y = (origHeight - p.Height) / 2
		}
	}

	// 创建新图像
	cutImg := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))

	// 复制像素
	for dy := 0; dy < p.Height; dy++ {
		for dx := 0; dx < p.Width; dx++ {
			cutImg.Set(dx, dy, img.At(x+dx, y+dy))
		}
	}

	return cutImg, nil
}

// NewCutProcessor 创建新的切割处理器（使用预定义位置）
func NewCutProcessor(width, height int, position CutPosition) *CutProcessor {
	// 验证位置参数
	validPositions := map[CutPosition]bool{
		CutPositionCenter: true,
		CutPositionTop:    true,
		CutPositionBottom: true,
		CutPositionLeft:   true,
		CutPositionRight:  true,
	}

	if !validPositions[position] {
		position = CutPositionCenter // 默认居中切割
	}

	return &CutProcessor{
		Width:           width,
		Height:          height,
		Position:        position,
		UseCustomRegion: false,
	}
}

// NewCutProcessorWithRegion 创建新的切割处理器（使用自定义区域）
func NewCutProcessorWithRegion(width, height, x, y int) *CutProcessor {
	return &CutProcessor{
		Width:           width,
		Height:          height,
		X:               x,
		Y:               y,
		UseCustomRegion: true,
	}
}

// NewSquareCutProcessor 创建正方形切割处理器
func NewSquareCutProcessor(size int, position CutPosition) *CutProcessor {
	return NewCutProcessor(size, size, position)
}
