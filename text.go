package vimage

import (
	"image"
	"image/color"
	"image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TextOptions 定义文本处理器的选项
type TextOptions struct {
	Text      string
	Position  image.Point
	Font      font.Face
	Color     color.Color
}

// DefaultTextOptions 默认文本选项
var DefaultTextOptions = TextOptions{
	Font:  basicfont.Face7x13,
	Color: color.Black,
}

// TextProcessor 实现文本处理器
type TextProcessor struct {
	Options TextOptions
}

// NewTextProcessor 创建新的文本处理器
func NewTextProcessor(opts TextOptions) *TextProcessor {
	if opts.Font == nil {
		opts.Font = DefaultTextOptions.Font
	}
	if opts.Color == nil {
		opts.Color = DefaultTextOptions.Color
	}
	return &TextProcessor{Options: opts}
}

// Process 实现ImageProcessor接口
func (p *TextProcessor) Process(img image.Image) (image.Image, error) {
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(p.Options.Color),
		Face: p.Options.Font,
		Dot:  fixed.P(p.Options.Position.X, p.Options.Position.Y),
	}
	drawer.DrawString(p.Options.Text)

	return dst, nil
}