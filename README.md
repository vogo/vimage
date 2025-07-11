# vimage

图像处理工具库，提供丰富的图像处理功能，包括调整大小、裁剪、马赛克、水印、噪点、验证码生成等。

## 安装

```bash
go get github.com/vogo/vimage
```

## 主要功能

- 图像缩放 (Zoom) - 提供精确的像素级缩放
- 图像切割 (Cut) - 从图像中切割指定区域
- 正方形裁剪 (Square)
- 圆形裁剪 (Circle)
- 马赛克处理 (Mosaic)
- 水印添加 (Watermark)
- 图像叠加 (Overlay)
- 噪点生成 (Noise)
- 验证码生成 (Captcha)
- 表格生成 (Table)

## 核心API用法示例

### 图像处理器框架

```go
// 创建处理器链
processors := []vimage.ImageProcessor{
    // 添加多个处理器
    processor1,
    processor2,
    // ...
}

// 处理图片
result, err := vimage.ProcessImage(imgData, processors, nil)
```

### 图像缩放 (Zoom)

```go
// 精确缩放（指定宽高）
zoomProcessor := vimage.NewZoomProcessor(width, height)

// 按比例缩放（例如：缩小到原来的50%）
zoomProcessor := vimage.NewZoomRatioProcessor(0.5)

// 按宽度缩放（高度按比例计算）
zoomProcessor := vimage.NewZoomWidthProcessor(300)

// 按高度缩放（宽度按比例计算）
zoomProcessor := vimage.NewZoomHeightProcessor(200)

// 按最大边缩放（保持比例）
zoomProcessor := vimage.NewZoomMaxProcessor(500)

// 按最小边缩放（保持比例）
zoomProcessor := vimage.NewZoomMinProcessor(300)

// 可选：设置缩放算法（默认为双线性插值）
zoomProcessor.WithScaler(draw.BiLinear) // 可选值: draw.NearestNeighbor, draw.ApproxBiLinear, draw.BiLinear, draw.CatmullRom

// 处理图像
zoomedImg, err := zoomProcessor.Process(srcImg)
```

### 图像切割 (Cut)

```go
// 使用预定义位置切割图像
// 位置可选: "center", "top", "bottom", "left", "right"
cutProcessor := vimage.NewCutProcessor(width, height, vimage.CutPositionCenter)

// 使用自定义区域切割图像
// x, y 是左上角坐标
cutProcessor := vimage.NewCutProcessorWithRegion(width, height, x, y)

// 创建正方形切割处理器（便捷方法）
cutProcessor := vimage.NewSquareCutProcessor(size, vimage.CutPositionCenter)

// 处理图像
cutImg, err := cutProcessor.Process(srcImg)
```



### 正方形裁剪

```go
// 创建正方形裁剪处理器，支持不同裁剪位置
// 位置可选: "center", "top", "bottom", "left", "right"
squareProcessor := vimage.NewSquareProcessor("center")

// 处理图像
squareImg, err := squareProcessor.Process(srcImg)
```

### 圆形裁剪

```go
// 创建圆形裁剪处理器
circleProcessor := &vimage.CircleProcessor{}

// 处理图像（注意：输入图像必须是正方形）
circleImg, err := circleProcessor.Process(squareImg)
```

### 马赛克处理

```go
// 定义马赛克区域
regions := []*vimage.MosaicRegion{
    {
        FromX: x1,
        FromY: y1,
        ToX:   x2,
        ToY:   y2,
    },
}

// 处理图像
result, err := vimage.MosaicImage(imgData, regions)

// 或使用更多选项
result, err := vimage.MosaicImageWithOptions(imgData, regions, 0.5, vimage.DirectionLeft)
```

### 水印添加

```go
// 创建水印处理器
watermarkProcessor := &vimage.WatermarkProcessor{
    Text:     "水印文本",
    FontSize: 24,
    Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
    Opacity:  0.7,
    Position: "bottom-right",
    Rotation: 30,
}

// 处理图像
result, err := watermarkProcessor.Process(srcImg)
```

### 图像叠加

```go
// 创建叠加处理器
overlayProcessor := &vimage.OverlayProcessor{
    OverlayImage: overlayImg,
    Position:     "center",
    Opacity:      0.8,
    Scale:        0.5,
}

// 处理图像
result, err := overlayProcessor.Process(srcImg)
```

### 验证码生成

```go
// 使用默认配置生成验证码
captchaText := "1234"
captchaImg, err := vimage.GenerateCaptcha(captchaText, nil)

// 使用自定义配置
config := &vimage.CaptchaConfig{
    Width:      160,
    Height:     60,
    NoiseLines: 8,
    NoiseDots:  100,
    BgColor:    color.RGBA{R: 240, G: 240, B: 240, A: 255},
    TextColor:  color.RGBA{R: 0, G: 0, B: 200, A: 255},
}
captchaImg, err := vimage.GenerateCaptcha(captchaText, config)
```

### 组合使用示例

```go
// 将图片裁剪为正方形并缩放
func SquareAndZoomImage(imgData []byte, position string, size int) ([]byte, error) {
    // 创建处理器链
    processors := []vimage.ImageProcessor{
        // 先裁剪为正方形
        vimage.NewSquareProcessor(position),
        // 再缩放
        vimage.NewZoomProcessor(size, size),
    }

    // 处理图片
    return vimage.ProcessImage(imgData, processors, nil)
}

// 将图片裁剪为正方形并应用圆形裁剪
func SquareAndCircleImage(imgData []byte, position string) ([]byte, error) {
    // 创建处理器链
    processors := []vimage.ImageProcessor{
        // 先裁剪为正方形
        vimage.NewSquareProcessor(position),
        // 再应用圆形裁剪
        &vimage.CircleProcessor{},
    }

    // 处理图片
    return vimage.ProcessImage(imgData, processors, nil)
}

// 按比例缩放图片并添加水印
func ZoomRatioAndWatermark(imgData []byte, ratio float64, watermarkText string) ([]byte, error) {
    // 创建处理器链
    processors := []vimage.ImageProcessor{
        // 先按比例缩放
        vimage.NewZoomRatioProcessor(ratio), // 使用新的缩放处理器
        // 再添加水印
        &vimage.WatermarkProcessor{
            Text:     watermarkText,
            FontSize: 24,
            Color:    image.White, // 使用预定义的白色
            Opacity:  0.7,
            Position: "bottom-right",
        },
    }

    // 处理图片
    return vimage.ProcessImage(imgData, processors, nil)
}

// 按最大边切割图片并裁剪为圆形
func CutAndCircle(imgData []byte, maxSize int) ([]byte, error) {
    // 创建处理器链
    processors := []vimage.ImageProcessor{
        // 先裁剪为正方形
        vimage.NewCutProcessor(maxSize, maxSize, vimage.CutPositionCenter), // 使用新的切割处理器
        // 最后裁剪为圆形
        &vimage.CircleProcessor{},
    }

    // 处理图片
    return vimage.ProcessImage(imgData, processors, nil)
}

// 先切割指定区域再缩放
func CutAndZoom(imgData []byte, cutWidth, cutHeight, x, y int, zoomRatio float64) ([]byte, error) {
    // 创建处理器链
    processors := []vimage.ImageProcessor{
        // 先切割指定区域
        vimage.NewCutProcessorWithRegion(cutWidth, cutHeight, x, y),
        // 再按比例缩放
        vimage.NewZoomRatioProcessor(zoomRatio),
    }

    // 处理图片
    return vimage.ProcessImage(imgData, processors, nil)
}
```

## 贡献代码

欢迎贡献代码，请遵循以下步骤：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建一个 Pull Request

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/yourusername/vimage.git
cd vimage

# 安装依赖
go mod download

# 安装开发工具
go install github.com/vogo/license-header-checker/cmd/license-header-checker@latest
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 代码规范

- 所有代码必须通过 `golangci-lint` 检查
- 所有新文件必须包含 Apache License 2.0 许可证头
- 所有函数和类型必须有适当的文档注释
- 测试覆盖率应尽可能高

### 构建和测试

```bash
# 格式化代码
make format

# 检查许可证头
make license-check

# 运行代码检查
make lint

# 运行测试
make test

# 完整构建流程（包括上述所有步骤）
make build
```

## 许可证

Apache License 2.0
