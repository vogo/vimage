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
	"math"

	"golang.org/x/image/draw"
)

// ZoomMode 定义缩放模式
type ZoomMode int

const (
	// ZoomModeExact 精确缩放到指定尺寸
	ZoomModeExact ZoomMode = iota
	// ZoomModeRatio 按比例缩放
	ZoomModeRatio
	// ZoomModeWidth 按宽度缩放，高度等比例调整
	ZoomModeWidth
	// ZoomModeHeight 按高度缩放，宽度等比例调整
	ZoomModeHeight
	// ZoomModeMax 按最大边缩放，保持比例
	ZoomModeMax
	// ZoomModeMin 按最小边缩放，保持比例
	ZoomModeMin
)

// ZoomProcessor 图像缩放处理器
// 对图像进行像素级缩放，不进行裁剪
type ZoomProcessor struct {
	// 目标宽度和高度
	Width  int
	Height int
	// 缩放比例 (0.0-1.0表示缩小，>1.0表示放大)
	Ratio float64
	// 缩放模式
	Mode ZoomMode
	// 缩放算法
	Scaler draw.Scaler
}

// Process 实现ImageProcessor接口
func (p *ZoomProcessor) Process(img image.Image) (image.Image, error) {
	// 获取原始图片尺寸
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 计算目标尺寸
	targetWidth, targetHeight := p.calculateTargetSize(origWidth, origHeight)

	// 验证目标尺寸
	if targetWidth <= 0 || targetHeight <= 0 {
		return nil, fmt.Errorf("无效的缩放尺寸: %dx%d", targetWidth, targetHeight)
	}

	// 创建目标图像
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// 使用指定的缩放算法
	scaler := p.Scaler
	if scaler == nil {
		// 默认使用双线性插值算法
		scaler = draw.BiLinear
	}

	// 执行缩放
	scaler.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst, nil
}

// calculateTargetSize 根据缩放模式计算目标尺寸
func (p *ZoomProcessor) calculateTargetSize(origWidth, origHeight int) (int, int) {
	switch p.Mode {
	case ZoomModeExact:
		// 精确缩放到指定尺寸
		return p.Width, p.Height

	case ZoomModeRatio:
		// 按比例缩放
		return int(math.Round(float64(origWidth) * p.Ratio)),
			int(math.Round(float64(origHeight) * p.Ratio))

	case ZoomModeWidth:
		// 按宽度缩放，高度等比例调整
		ratio := float64(p.Width) / float64(origWidth)
		return p.Width, int(math.Round(float64(origHeight) * ratio))

	case ZoomModeHeight:
		// 按高度缩放，宽度等比例调整
		ratio := float64(p.Height) / float64(origHeight)
		return int(math.Round(float64(origWidth) * ratio)), p.Height

	case ZoomModeMax:
		// 按最大边缩放，保持比例
		var ratio float64
		if origWidth >= origHeight {
			// 宽度是最大边
			ratio = float64(p.Width) / float64(origWidth)
		} else {
			// 高度是最大边
			ratio = float64(p.Height) / float64(origHeight)
		}
		return int(math.Round(float64(origWidth) * ratio)),
			int(math.Round(float64(origHeight) * ratio))

	case ZoomModeMin:
		// 按最小边缩放，保持比例
		var ratio float64
		if origWidth <= origHeight {
			// 宽度是最小边
			ratio = float64(p.Width) / float64(origWidth)
		} else {
			// 高度是最小边
			ratio = float64(p.Height) / float64(origHeight)
		}
		return int(math.Round(float64(origWidth) * ratio)),
			int(math.Round(float64(origHeight) * ratio))

	default:
		// 默认精确缩放
		return p.Width, p.Height
	}
}

// NewZoomProcessor 创建新的精确缩放处理器
func NewZoomProcessor(width, height int) *ZoomProcessor {
	return &ZoomProcessor{
		Width:  width,
		Height: height,
		Mode:   ZoomModeExact,
	}
}

// NewZoomRatioProcessor 创建新的按比例缩放处理器
func NewZoomRatioProcessor(ratio float64) *ZoomProcessor {
	return &ZoomProcessor{
		Ratio: ratio,
		Mode:  ZoomModeRatio,
	}
}

// NewZoomWidthProcessor 创建新的按宽度缩放处理器
func NewZoomWidthProcessor(width int) *ZoomProcessor {
	return &ZoomProcessor{
		Width: width,
		Mode:  ZoomModeWidth,
	}
}

// NewZoomHeightProcessor 创建新的按高度缩放处理器
func NewZoomHeightProcessor(height int) *ZoomProcessor {
	return &ZoomProcessor{
		Height: height,
		Mode:   ZoomModeHeight,
	}
}

// NewZoomMaxProcessor 创建新的按最大边缩放处理器
func NewZoomMaxProcessor(size int) *ZoomProcessor {
	return &ZoomProcessor{
		Width:  size,
		Height: size,
		Mode:   ZoomModeMax,
	}
}

// NewZoomMinProcessor 创建新的按最小边缩放处理器
func NewZoomMinProcessor(size int) *ZoomProcessor {
	return &ZoomProcessor{
		Width:  size,
		Height: size,
		Mode:   ZoomModeMin,
	}
}

// WithScaler 设置缩放算法
func (p *ZoomProcessor) WithScaler(scaler draw.Scaler) *ZoomProcessor {
	p.Scaler = scaler
	return p
}
