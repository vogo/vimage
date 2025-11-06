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
	"image/jpeg"
	"image/png"

	"github.com/fogleman/gg"
)

// Processor 定义新的处理器接口，支持选项参数
type Processor interface {
	Process(img image.Image) (image.Image, error)
}

type ImageProcessContext struct {
	dc     *gg.Context
	Width  int
	Height int
}

func (ctx *ImageProcessContext) DC() *gg.Context {
	return ctx.dc
}

func NewImageProcessContext(img image.Image) *ImageProcessContext {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	return &ImageProcessContext{
		dc:     dc,
		Width:  width,
		Height: height,
	}
}

type ContextProcessor interface {
	ContextProcess(ctx *ImageProcessContext) error
}

// ProcessorOptions 处理器选项
type ProcessorOptions struct {
	// 可以添加通用选项
	Quality int // JPEG压缩质量 (1-100)
}

// DefaultProcessorOptions 默认处理器选项
var DefaultProcessorOptions = ProcessorOptions{
	Quality: 90,
}

// ProcessImage 使用处理器链处理图片
// imgData: 原始图片字节数据
// processors: 处理器链
// options: 处理选项
// 返回: 处理后的图片字节数据和错误信息
func ProcessImage(imgData []byte, processors []Processor, options *ProcessorOptions) ([]byte, error) {
	// 使用默认选项
	if options == nil {
		options = &DefaultProcessorOptions
	}

	// 解码图片
	srcImg, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	// 应用处理器链
	currentImg, err := Process(srcImg, processors)
	if err != nil {
		return nil, err
	}

	// 编码图片
	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, currentImg, &jpeg.Options{Quality: options.Quality})
	case "png":
		err = png.Encode(buf, currentImg)
	default:
		// 默认使用PNG格式
		err = png.Encode(buf, currentImg)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Process 循环处理图片
func Process(img image.Image, processors []Processor) (image.Image, error) {
	var err error
	currentImg := img
	for _, processor := range processors {
		currentImg, err = processor.Process(currentImg)
		if err != nil {
			return nil, err
		}
	}

	return currentImg, nil
}

// ContextProcess 上下文处理
func ContextProcess(img image.Image, processors []ContextProcessor) (image.Image, error) {
	ctx := NewImageProcessContext(img)

	for _, processor := range processors {
		if err := processor.ContextProcess(ctx); err != nil {
			return nil, err
		}
	}

	return ctx.dc.Image(), nil
}

type EmptyProcessor struct{}

func (p *EmptyProcessor) Process(img image.Image) (image.Image, error) {
	return img, nil
}
