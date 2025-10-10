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

var (
	notoSansSCWghtFont   *truetype.Font
	harmonyOSSansSCBlack *truetype.Font

	loadFontMutex = sync.Mutex{}
)

func LoadNotoSansSCVariableFontWght() *truetype.Font {
	return LoadFont(&notoSansSCWghtFont, "/opt/NotoSansSC_wght.ttf", "https://raw.githubusercontent.com/google/fonts/refs/heads/main/ofl/notosanssc/NotoSansSC%5Bwght%5D.ttf")
}

func LoadHarmonyOSSansSCBlack() *truetype.Font {
	return LoadFont(&harmonyOSSansSCBlack, "/opt/HarmonyOS_Sans_SC_Black.ttf", "")
}

func LoadFont(font **truetype.Font, localPath, downloadUrl string) *truetype.Font {
	loadFontMutex.Lock()
	defer loadFontMutex.Unlock()

	if *font != nil {
		return *font
	}

	// check file exists
	if _, err := os.Stat(localPath); err == nil {
		fontBytes, err := os.ReadFile(localPath)
		if err == nil {
			font, err := truetype.Parse(fontBytes)
			if err == nil {
				fmt.Printf("load font from local file: %s\n", localPath)
				notoSansSCWghtFont = font
				return font
			}
		}
	}

	if downloadUrl == "" {
		panic("font not found " + localPath)
	}

	fmt.Printf("download font from: %s\n", downloadUrl)

	// 创建一个带有60秒超时的HTTP客户端
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(downloadUrl)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	fontBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	_ = os.WriteFile(localPath, fontBytes, 0o644)

	fontObj, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	*font = fontObj

	return fontObj
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
