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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCutCircleProcessor(t *testing.T) {
	// 创建一个 200x200 的 RGBA 图片 (正方形)
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// 填充一些颜色以便测试
	for y := range 200 {
		for x := range 200 {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	circleProcessor := NewCutCircleProcessor()
	newImg, err := circleProcessor.Process(img)
	require.NoError(t, err)

	// 检查结果图片是正方形
	bounds := newImg.Bounds()
	assert.Equal(t, bounds.Dx(), bounds.Dy(), "Output image should be square")

	// 检查图片尺寸与输入一致
	assert.Equal(t, 200, bounds.Dx(), "Output width should match input")
	assert.Equal(t, 200, bounds.Dy(), "Output height should match input")

	// 检查中心点的透明度 (应该是不透明的)
	centerColor := newImg.At(100, 100).(color.RGBA)
	assert.Equal(t, uint8(255), centerColor.A, "Center should be opaque")

	// 检查角落的透明度 (应该是透明的)
	cornerColor := newImg.At(0, 0).(color.RGBA)
	assert.Equal(t, uint8(0), cornerColor.A, "Corner should be transparent")
}

func TestCutCircleProcessor_NotSquare(t *testing.T) {
	// 创建一个 300x200 的非正方形 RGBA 图片
	img := image.NewRGBA(image.Rect(0, 0, 300, 200))

	// 填充一些颜色以便测试
	for y := range 200 {
		for x := range 300 {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	circleProcessor := NewCutCircleProcessor()
	_, err := circleProcessor.Process(img)

	// CutCircleProcessor 要求输入必须是正方形，非正方形应该返回错误
	require.Error(t, err, "Non-square image should return an error")
	assert.Contains(t, err.Error(), "square", "Error message should mention square requirement")
}
