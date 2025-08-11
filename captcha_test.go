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
	"image/color"
	"os"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/basicfont"
)

// TestGenCaptchaImage 测试验证码图片生成
func TestGenCaptchaImage(t *testing.T) {
	testCases := []struct {
		name    string
		captcha string
	}{
		{"数字验证码", "1234"},
		{"字母验证码", "ABCD"},
		{"混合验证码", "A1B2"},
		{"长验证码", "123456"},
		{"短验证码", "AB"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := GenCaptchaImage(tc.captcha)
			if err != nil {
				t.Fatalf("生成验证码图片失败: %v", err)
			}

			if buf == nil {
				t.Fatal("返回的缓冲区为空")
			}

			if buf.Len() == 0 {
				t.Fatal("生成的图片数据为空")
			}

			// 可选：保存图片到文件进行手动检查
			if os.Getenv("SAVE_TEST_IMAGES") == "true" {
				filename := "build/test_captcha_" + tc.captcha + ".png"
				if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
					t.Logf("保存测试图片失败: %v", err)
				} else {
					t.Logf("测试图片已保存: %s", filename)
				}
			}
		})
	}
}

// TestGenCaptchaImageWithConfig 测试自定义配置的验证码生成
func TestGenCaptchaImageWithConfig(t *testing.T) {
	// 自定义配置
	customConfig := &CaptchaConfig{
		Width:            200,
		Height:           60,
		BgColor:          color.RGBA{R: 240, G: 240, B: 240, A: 255}, // 浅灰背景
		TextColor:        color.RGBA{R: 50, G: 50, B: 200, A: 255},   // 蓝色文字
		NoiseLines:       8,
		NoiseDots:        80,
		Face:             basicfont.Face7x13,
		CharSpacing:      45, // 大字体的字符间距
		CharWidth:        40, // 大字体的字符宽度
		CharYOffsetRange: 12, // 大字体的垂直偏移范围
		CharXOffsetRange: 8,  // 大字体的水平偏移范围
	}

	font, err := LoadNotoSansSCVariableFontWght()
	if err != nil {
		t.Fatalf("加载字体失败: %v", err)
	}
	customConfig.Face = truetype.NewFace(font, &truetype.Options{
		Size: 32,
	})

	buf, err := GenCaptchaImageWithConfig("TEST", customConfig)
	if err != nil {
		t.Fatalf("使用自定义配置生成验证码失败: %v", err)
	}

	if buf == nil || buf.Len() == 0 {
		t.Fatal("生成的图片数据无效")
	}

	// 可选：保存自定义配置的图片
	if os.Getenv("SAVE_TEST_IMAGES") == "true" {
		if err := os.WriteFile("build/test_captcha_custom.png", buf.Bytes(), 0o644); err != nil {
			t.Logf("保存自定义配置测试图片失败: %v", err)
		} else {
			t.Log("自定义配置测试图片已保存: test_captcha_custom.png")
		}
	}
}

// BenchmarkGenCaptchaImage 性能测试
func BenchmarkGenCaptchaImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenCaptchaImage("1234")
		if err != nil {
			b.Fatalf("生成验证码失败: %v", err)
		}
	}
}
