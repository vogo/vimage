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
	_ "embed"
	"fmt"
	"image/color"
	"image/png"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/basicfont"
)

// GenMultipleColumnsTableImage 根据数据生成列式表格图片
// headers 为标题，作为最左侧的一列显示
// data 为数据内容，每一条数据作为一列展示
// 返回PNG格式的图片数据
func GenMultipleColumnsTableImage(font *truetype.Font, headers []string, data [][]string) (*bytes.Buffer, error) {
	// 输入验证
	if len(headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	// 验证每条数据的长度是否与表头一致
	for i, row := range data {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("data %d has %d items, expected %d", i, len(row), len(headers))
		}
	}

	// 表格参数
	defaultColWidth := 120.0 // 默认列宽
	headerColWidth := 150.0  // header列宽度（稍宽一些）
	rowHeight := 40.0        // 行高

	// 创建图片
	// 宽度 = header列宽 + 数据列数 * 默认列宽
	imgWidth := int(headerColWidth + float64(len(data))*defaultColWidth)
	// 高度 = header行数 * 行高
	imgHeight := int(float64(len(headers)) * rowHeight)
	dc := gg.NewContext(imgWidth, imgHeight)

	// 设置背景色
	dc.SetColor(color.White)
	dc.Clear()

	// 加载字体
	if font != nil {
		face := truetype.NewFace(font, &truetype.Options{Size: 14})
		dc.SetFontFace(face)
	} else {
		dc.SetFontFace(basicfont.Face7x13)
	}

	// 颜色定义
	headerBg := color.RGBA{R: 50, G: 100, B: 200, A: 255}     // header列背景色
	headerFg := color.White                                   // header列文字颜色
	dataBg := color.RGBA{R: 240, G: 245, B: 255, A: 255}      // 数据列背景色
	dataAltBg := color.White                                  // 数据列交替背景色
	borderColor := color.RGBA{R: 200, G: 200, B: 200, A: 255} // 边框颜色

	// 绘制header列（最左侧）
	dc.SetColor(headerBg)
	dc.DrawRectangle(0, 0, headerColWidth, float64(imgHeight))
	dc.Fill()

	// 绘制header列文本
	dc.SetColor(headerFg)
	for i, header := range headers {
		yPos := float64(i)*rowHeight + rowHeight/2
		dc.DrawStringAnchored(header, headerColWidth/2, yPos, 0.5, 0.5)
	}

	// 绘制数据列
	for colIdx, dataCol := range data {
		// 交替列背景色
		if colIdx%2 == 0 {
			dc.SetColor(dataBg)
		} else {
			dc.SetColor(dataAltBg)
		}

		xPos := headerColWidth + float64(colIdx)*defaultColWidth
		dc.DrawRectangle(xPos, 0, defaultColWidth, float64(imgHeight))
		dc.Fill()

		// 绘制数据列文本
		dc.SetColor(color.Black)
		for rowIdx, cell := range dataCol {
			yPos := float64(rowIdx)*rowHeight + rowHeight/2
			dc.DrawStringAnchored(cell, xPos+defaultColWidth/2, yPos, 0.5, 0.5)
		}
	}

	// 绘制行分隔线
	dc.SetColor(borderColor)
	dc.SetLineWidth(1)
	for i := 1; i < len(headers); i++ {
		yPos := float64(i) * rowHeight
		dc.DrawLine(0, yPos, float64(imgWidth), yPos)
		dc.Stroke()
	}

	// 绘制列分隔线
	// header列右侧分隔线
	dc.DrawLine(headerColWidth, 0, headerColWidth, float64(imgHeight))
	dc.Stroke()

	// 数据列分隔线
	for i := 1; i < len(data); i++ {
		xPos := headerColWidth + float64(i)*defaultColWidth
		dc.DrawLine(xPos, 0, xPos, float64(imgHeight))
		dc.Stroke()
	}

	// 绘制外边框
	dc.SetLineWidth(2)
	dc.DrawRectangle(0, 0, float64(imgWidth), float64(imgHeight))
	dc.Stroke()

	buf := new(bytes.Buffer)
	err := png.Encode(buf, dc.Image())
	if err != nil {
		return nil, fmt.Errorf("图片编码失败: %w", err)
	}

	return buf, nil
}

// GenMultipleRowsTableImage 根据数据生成表格图片
// headers 为标题， data 为数据内容
// widths 为每列宽度，如果为空则使用默认宽度
// 返回PNG格式的图片数据
func GenMultipleRowsTableImage(font *truetype.Font, headers []string, data [][]string, widths []float64) (*bytes.Buffer, error) {
	// 输入验证
	if len(headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	// 验证每行的列数是否与表头一致
	for i, row := range data {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("row %d has %d columns, expected %d", i, len(row), len(headers))
		}
	}

	// 表格参数
	defaultColWidth := 120.0 // 默认列宽
	rowHeight := 40.0        // 行高
	headerHeight := 50.0     // 表头高度
	padding := 10.0          // 内边距

	// 设置列宽
	colWidths := make([]float64, len(headers))
	if len(widths) == len(headers) {
		// 使用传入的列宽
		copy(colWidths, widths)
	} else {
		// 使用默认列宽
		for i := range colWidths {
			colWidths[i] = defaultColWidth
		}
	}

	// 计算总宽度
	totalWidth := 0.0
	for _, width := range colWidths {
		totalWidth += width
	}

	// 创建图片
	imgWidth := int(totalWidth)
	imgHeight := int(headerHeight + float64(len(data))*rowHeight)
	dc := gg.NewContext(imgWidth, imgHeight)

	// 设置背景色
	dc.SetColor(color.White)
	dc.Clear()

	// 加载字体

	if font != nil {
		face := truetype.NewFace(font, &truetype.Options{Size: 14})
		dc.SetFontFace(face)
	} else {
		dc.SetFontFace(basicfont.Face7x13)
	}

	// 绘制表头
	headerBg := color.RGBA{R: 50, G: 100, B: 200, A: 255}
	headerFg := color.White
	dc.SetColor(headerBg)
	dc.DrawRectangle(0, 0, float64(imgWidth), headerHeight)
	dc.Fill()

	y := headerHeight/2 - 7 // 垂直居中偏移
	x := padding
	for i, header := range headers {
		dc.SetColor(headerFg)
		dc.DrawStringAnchored(header, x+colWidths[i]/2, y, 0.5, 0.5)
		x += colWidths[i]
	}

	// 绘制表格行
	rowBg := color.RGBA{R: 240, G: 245, B: 255, A: 255}
	rowAltBg := color.White
	borderColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}

	for rowIdx, row := range data {
		// 交替行背景色
		if rowIdx%2 == 0 {
			dc.SetColor(rowBg)
		} else {
			dc.SetColor(rowAltBg)
		}

		yPos := headerHeight + float64(rowIdx)*rowHeight
		dc.DrawRectangle(0, yPos, float64(imgWidth), rowHeight)
		dc.Fill()

		// 绘制单元格文本
		dc.SetColor(color.Black)
		x = padding
		for colIdx, cell := range row {
			dc.SetColor(color.Black)

			dc.DrawStringAnchored(
				cell,
				x+colWidths[colIdx]/2,
				yPos+rowHeight/2,
				0.5,
				0.5,
			)
			x += colWidths[colIdx]
		}

		// 绘制行边框
		dc.SetColor(borderColor)
		dc.SetLineWidth(1)
		dc.DrawLine(0, yPos, float64(imgWidth), yPos)
		dc.Stroke()
	}

	// 绘制列分隔线
	dc.SetColor(borderColor)
	x = 0
	for _, width := range colWidths[:len(colWidths)-1] {
		x += width
		dc.DrawLine(x, headerHeight, x, float64(imgHeight))
		dc.Stroke()
	}

	// 绘制外边框
	dc.SetLineWidth(2)
	dc.DrawRectangle(0, 0, float64(imgWidth), float64(imgHeight))
	dc.Stroke()

	buf := new(bytes.Buffer)
	err := png.Encode(buf, dc.Image())
	if err != nil {
		log.Fatal("图片编码失败:", err)
	}

	return buf, nil
}
