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

func TestDrawCircleProcessor(t *testing.T) {
	// 创建一个 200x200 的 RGBA 图片
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// 初始化 DrawCircleProcessor (不填充)
	// 在中心 (100, 100) 绘制一个半径为 50 的红色圆
	processor := NewDrawCircleProcessor(100, 100, 50, color.RGBA{R: 255, A: 255}, false)

	// 处理图片
	processedImg, err := processor.Process(img)
	require.NoError(t, err)

	// 检查圆上某个点的颜色
	// 例如，检查点 (150, 100)，它应该在圆的边缘
	// 由于抗锯齿的存在，颜色可能不是纯红色，但红色分量应该很高
	x, y := 150, 100
	c := processedImg.At(x, y).(color.RGBA)

	// 断言红色分量大于0
	assert.True(t, c.R > 0, "The R component of the color should be greater than 0")

	// 检查一个不在圆内的点，例如 (0, 0)
	c = processedImg.At(0, 0).(color.RGBA)
	assert.Equal(t, uint8(0), c.R, "The R component of the color should be 0 for a point outside the circle")
}

func TestDrawCircleProcessorFilled(t *testing.T) {
	// 创建一个 200x200 的 RGBA 图片
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// 初始化 DrawCircleProcessor (填充)
	// 在中心 (100, 100) 绘制一个半径为 50 的填充红色圆
	processor := NewDrawCircleProcessor(100, 100, 50, color.RGBA{R: 255, A: 255}, true)

	// 处理图片
	processedImg, err := processor.Process(img)
	require.NoError(t, err)

	// 检查圆内的点，例如 (100, 100)，应该是红色
	c := processedImg.At(100, 100).(color.RGBA)
	assert.True(t, c.R > 0, "The R component should be greater than 0 for a point inside the filled circle")

	// 检查圆内的另一个点 (120, 100)
	c = processedImg.At(120, 100).(color.RGBA)
	assert.True(t, c.R > 0, "The R component should be greater than 0 for a point inside the filled circle")

	// 检查一个不在圆内的点，例如 (0, 0)
	c = processedImg.At(0, 0).(color.RGBA)
	assert.Equal(t, uint8(0), c.R, "The R component should be 0 for a point outside the circle")
}
