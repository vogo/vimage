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
	"image/png"
	"os"
	"testing"
)

// 创建测试图片
func createTestImageForProcessor(width, height int) []byte {
	// 创建一个彩色渐变图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 创建彩色渐变效果
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// 编码为PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

// TestProcessImage 测试图片处理器框架
func TestProcessImage(t *testing.T) {
	// 创建测试图片
	testImg := createTestImageForProcessor(200, 200)

	// 创建处理器链
	processors := []ImageProcessor{
		// 添加马赛克处理器
		NewMosaicProcessor([]*MosaicRegion{
			{
				FromX: 50,
				FromY: 50,
				ToX:   150,
				ToY:   150,
			},
		}, 1.0, DirectionLeft),

		// 添加水印处理器
		NewWatermarkProcessor("测试水印", 20, color.RGBA{R: 255, G: 255, B: 255, A: 255}, 0.7, "bottom-right", 0),
	}

	// 处理图片
	result, err := ProcessImage(testImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("处理后的图片数据为空")
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_processor.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_processor.png")
	}
}

// TestMultipleProcessors 测试多个处理器链式处理
func TestMultipleProcessors(t *testing.T) {
	// 创建测试图片
	testImg := createTestImageForProcessor(300, 200)

	// 创建处理器链
	processors := []ImageProcessor{
		// 添加马赛克处理器
		NewMosaicProcessor([]*MosaicRegion{
			{
				FromX: 20,
				FromY: 20,
				ToX:   80,
				ToY:   80,
			},
			{
				FromX: 120,
				FromY: 120,
				ToX:   180,
				ToY:   180,
			},
		}, 1.0, DirectionLeft),

		// 添加噪点处理器
		NewNoiseProcessor(
			5,  // 干扰线数量
			50, // 干扰点数量
			color.RGBA{R: 200, G: 200, B: 200, A: 150}, // 线条颜色
			color.RGBA{R: 200, G: 200, B: 200, A: 150}, // 点颜色
		),

		// 添加水印处理器
		NewWatermarkProcessor(
			"多重处理测试",
			24,
			color.RGBA{R: 255, G: 100, B: 100, A: 255},
			0.6,
			"center",
			-15, // 旋转角度
		),
	}

	// 处理图片
	result, err := ProcessImage(testImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("处理后的图片数据为空")
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_multiple_processors.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_multiple_processors.png")
	}
}

// TestResizeProcessor 测试调整大小处理器
func TestResizeProcessor(t *testing.T) {
	// 创建测试图片
	testImg := createTestImageForProcessor(400, 300)

	// 创建处理器链
	processors := []ImageProcessor{
		// 添加调整大小处理器
		NewResizeProcessor(200, 150),
	}

	// 处理图片
	result, err := ProcessImage(testImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	// 验证返回的是有效的图片数据
	decodedImg, _, err := image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}

	// 验证图片尺寸
	bounds := decodedImg.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 150 {
		t.Fatalf("调整大小失败，期望尺寸 200x150，实际尺寸 %dx%d", bounds.Dx(), bounds.Dy())
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_resize_processor.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_resize_processor.png")
	}
}

// TestSquareProcessorInChain 测试正方形裁剪处理器在处理链中的使用
func TestSquareProcessorInChain(t *testing.T) {
	// 创建测试图片（非正方形）
	testImg := createTestImageForProcessor(400, 300)

	// 创建处理器链
	processors := []ImageProcessor{
		// 添加正方形裁剪处理器
		NewSquareProcessor("center"),
	}

	// 处理图片
	result, err := ProcessImage(testImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	// 验证返回的是有效的图片数据
	decodedImg, _, err := image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}

	// 验证图片是正方形
	bounds := decodedImg.Bounds()
	if bounds.Dx() != bounds.Dy() {
		t.Fatalf("裁剪为正方形失败，期望正方形，实际尺寸 %dx%d", bounds.Dx(), bounds.Dy())
	}

	// 验证尺寸是较小的一边（300）
	if bounds.Dx() != 300 {
		t.Fatalf("裁剪尺寸错误，期望尺寸 300x300，实际尺寸 %dx%d", bounds.Dx(), bounds.Dy())
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_square_processor.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_square_processor.png")
	}
}

// BenchmarkProcessImage 性能测试
func BenchmarkProcessImage(b *testing.B) {
	testImg := createTestImageForProcessor(200, 200)

	processors := []ImageProcessor{
		NewMosaicProcessor([]*MosaicRegion{
			{
				FromX: 50,
				FromY: 50,
				ToX:   150,
				ToY:   150,
			},
		}, 1.0, DirectionLeft),
		NewWatermarkProcessor("性能测试", 20, color.RGBA{R: 255, G: 255, B: 255, A: 255}, 0.7, "bottom-right", 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ProcessImage(testImg, processors, nil)
		if err != nil {
			b.Fatalf("图片处理失败: %v", err)
		}
	}
}
