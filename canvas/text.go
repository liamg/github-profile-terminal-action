package canvas

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (c *Canvas) MeasureText(input string, face font.Face) (int, int) {
	var width int
	fontHeight := int(face.Metrics().Height >> 6)
	height := fontHeight
	for _, r := range input {
		if r == '\n' {
			height += fontHeight
		}
		if advance, ok := face.GlyphAdvance(r); ok {
			width += int(advance >> 6)
		}
	}
	return width, height
}

func (c *Canvas) DrawText(p image.Point, col color.RGBA, face font.Face, text string) (int, int) {
	w, h := c.MeasureText(text, face)
	point := fixed.Point26_6{
		X: fixed.I(p.X),
		Y: fixed.I(p.Y + h),
	}
	d := &font.Drawer{
		Dst:  c.img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
	return w, h
}
