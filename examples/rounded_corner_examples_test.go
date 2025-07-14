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

package examples

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/vogo/vimage"
)

// ExampleRoundedCornerProcessor 展示如何使用圆角处理器
func TestRoundedCornerProcessor(t *testing.T) {
	// 创建一个测试图像 (200x200 的蓝色方块)
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255}) // 蓝色
		}
	}

	// 创建圆角处理器，设置圆角半径为30像素
	processor := vimage.NewRoundedCornerProcessor(30)

	// 处理图像
	result, err := processor.Process(img)
	if err != nil {
		panic(err)
	}

	// 将结果保存为PNG文件
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, result); err != nil {
		panic(err)
	}

	// 输出到文件
	if err := os.WriteFile("/tmp/rounded_corner_example.png", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	// 打印输出信息
	fmt.Println("圆角处理器示例已保存到/tmp/rounded_corner_example.png")
}
