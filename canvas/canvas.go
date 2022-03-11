package canvas

import (
	"image"
	"image/color"
	"image/draw"
)

type Canvas struct {
	width, height int
	img           *image.RGBA
}

func New(width, height int) *Canvas {
	return &Canvas{
		width:  width,
		height: height,
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (c *Canvas) Image() image.Image {
	return c.img
}

func (c *Canvas) Clone() *Canvas {
	clone := New(c.width, c.height)
	draw.Draw(clone.img, c.img.Bounds(), c.img, image.Point{}, draw.Src)
	return clone
}

func (c *Canvas) Fill(col color.RGBA) {
	c.Rect(0, 0, c.width, c.height, col)
}

func (c *Canvas) Rect(x1, y1, x2, y2 int, col color.RGBA) {
	rect := image.Rect(x1, y1, x2, y2)
	draw.Draw(c.img, rect, &image.Uniform{C: col}, image.Point{}, draw.Src)
}

func (c *Canvas) DrawImage(p image.Point, src image.Image) {
	draw.Draw(c.img, c.img.Bounds(), src, p, draw.Src)
}
