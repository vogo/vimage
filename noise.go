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
	"math/rand"
)

// NoiseProcessor 噪点处理器
type NoiseProcessor struct {
	NoiseLines int        // 干扰线数量
	NoiseDots  int        // 干扰点数量
	LineColor  color.RGBA // 线条颜色
	DotColor   color.RGBA // 点颜色
}

// Process 实现Processor接口
func (p *NoiseProcessor) Process(img image.Image) (image.Image, error) {
	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建新的RGBA图像
	dstImg := image.NewRGBA(bounds)

	// 复制原图像到新图像
	draw.Draw(dstImg, bounds, img, bounds.Min, draw.Src)

	// 添加干扰线
	for i := 0; i < p.NoiseLines; i++ {
		// 随机起点和终点
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)

		// 绘制线条
		DrawLine(dstImg, x1, y1, x2, y2, p.LineColor)
	}

	// 添加干扰点
	for i := 0; i < p.NoiseDots; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		dstImg.Set(x, y, p.DotColor)
	}

	return dstImg, nil
}

// NewNoiseProcessor 创建新的噪点处理器
func NewNoiseProcessor(noiseLines, noiseDots int, lineColor, dotColor color.RGBA) *NoiseProcessor {
	return &NoiseProcessor{
		NoiseLines: noiseLines,
		NoiseDots:  noiseDots,
		LineColor:  lineColor,
		DotColor:   dotColor,
	}
}
