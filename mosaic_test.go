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
	"testing"
)

// createTestImage 创建一个测试用的彩色图片
func createTestImage(width, height int) []byte {
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
		panic(err) // 在测试辅助函数中，可以使用panic处理错误
	}
	return buf.Bytes()
}

// TestMosaicImage 测试马赛克功能
func TestMosaicImage(t *testing.T) {
	// 创建测试图片
	testImg := createTestImage(100, 100)

	// 测试正常情况 - 使用向后兼容函数
	result, err := MosaicImageSingle(testImg, 20, 20, 80, 80)
	if err != nil {
		t.Fatalf("马赛克处理失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("处理后的图片数据为空")
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}
}

// TestMosaicImageInvalidCoordinates 测试无效坐标
func TestMosaicImageInvalidCoordinates(t *testing.T) {
	testImg := createTestImage(100, 100)

	// 测试无效坐标（fromX >= toX）
	regions := []*MosaicRegion{
		{
			FromX: 80,
			FromY: 20,
			ToX:   20,
			ToY:   80,
		},
	}
	_, err := MosaicImage(testImg, regions)
	// 现在我们不会返回错误，而是跳过无效区域
	if err != nil {
		t.Fatalf("应该跳过无效区域而不是返回错误: %v", err)
	}

	// 测试无效坐标（fromY >= toY）
	regions = []*MosaicRegion{
		{
			FromX: 20,
			FromY: 80,
			ToX:   80,
			ToY:   20,
		},
	}
	result, err := MosaicImage(testImg, regions)
	// 现在我们不会返回错误，而是跳过无效区域
	if err != nil {
		t.Fatalf("应该跳过无效区域而不是返回错误: %v", err)
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}
}

// TestMosaicImageBoundaryCorrection 测试边界自动修正
func TestMosaicImageBoundaryCorrection(t *testing.T) {
	testImg := createTestImage(100, 100)

	// 测试超出边界的坐标会被自动修正
	regions := []*MosaicRegion{
		{
			FromX: -10,
			FromY: -10,
			ToX:   150,
			ToY:   150,
		},
	}
	result, err := MosaicImage(testImg, regions)
	if err != nil {
		t.Fatalf("边界修正测试失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("处理后的图片数据为空")
	}
}

// TestMosaicImageInvalidData 测试无效图片数据
func TestMosaicImageInvalidData(t *testing.T) {
	invalidData := []byte("这不是图片数据")

	_, err := MosaicImageSingle(invalidData, 0, 0, 10, 10)
	if err == nil {
		t.Fatal("应该返回解码错误")
	}
}

// TestMosaicImageMultipleRegions 测试多个马赛克区域
func TestMosaicImageMultipleRegions(t *testing.T) {
	testImg := createTestImage(200, 200)

	// 定义多个马赛克区域
	regions := []*MosaicRegion{
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
		{
			FromX: 20,
			FromY: 120,
			ToX:   80,
			ToY:   180,
		},
	}

	result, err := MosaicImage(testImg, regions)
	if err != nil {
		t.Fatalf("多区域马赛克处理失败: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("处理后的图片数据为空")
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}
}

// BenchmarkMosaicImage 性能测试
func BenchmarkMosaicImage(b *testing.B) {
	testImg := createTestImage(200, 200)

	regions := []*MosaicRegion{
		{
			FromX: 50,
			FromY: 50,
			ToX:   150,
			ToY:   150,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MosaicImage(testImg, regions)
		if err != nil {
			b.Fatalf("马赛克处理失败: %v", err)
		}
	}
}

// BenchmarkMosaicImageMultipleRegions 多区域性能测试
func BenchmarkMosaicImageMultipleRegions(b *testing.B) {
	testImg := createTestImage(200, 200)

	regions := []*MosaicRegion{
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
		{
			FromX: 20,
			FromY: 120,
			ToX:   80,
			ToY:   180,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MosaicImage(testImg, regions)
		if err != nil {
			b.Fatalf("马赛克处理失败: %v", err)
		}
	}
}
