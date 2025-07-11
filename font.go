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
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
)

var notoSansSCWghtFont *truetype.Font
var notoSansSCWghtFontMutex = sync.Mutex{}

func LoadNotoSansSCVariableFontWght() (*truetype.Font, error) {
	notoSansSCWghtFontMutex.Lock()
	defer notoSansSCWghtFontMutex.Unlock()

	if notoSansSCWghtFont != nil {
		return notoSansSCWghtFont, nil
	}

	localPath := "/tmp/NotoSansSC_wght.ttf"

	// check file exists
	if _, err := os.Stat(localPath); err == nil {
		fontBytes, err := os.ReadFile(localPath)
		if err == nil {
			font, err := truetype.Parse(fontBytes)
			if err == nil {
				fmt.Printf("load font from local file: %s\n", localPath)
				return font, nil
			}
		}
	}

	url := "https://raw.githubusercontent.com/google/fonts/refs/heads/main/ofl/notosanssc/NotoSansSC%5Bwght%5D.ttf"
	fmt.Printf("download font from: %s\n", url)

	// 创建一个带有60秒超时的HTTP客户端
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fontBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = os.WriteFile(localPath, fontBytes, 0o644)

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	notoSansSCWghtFont = font

	return font, nil
}

var defaultFont *truetype.Font

func SetDefaultFont(font *truetype.Font) {
	defaultFont = font
}

func GetDefaultFont() (*truetype.Font, error) {
	if defaultFont == nil {
		return nil, fmt.Errorf("default font not set")
	}

	return defaultFont, nil
}

func requireFont() *truetype.Font {
	if defaultFont != nil {
		return defaultFont
	}
	font, err := LoadNotoSansSCVariableFontWght()
	if err != nil {
		panic(err)
	}
	return font
}
