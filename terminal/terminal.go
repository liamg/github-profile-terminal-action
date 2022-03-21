package terminal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/liamg/github-profile-magic-action/canvas"
	"github.com/liamg/github-profile-magic-action/theme"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Speed int

const (
	Slow     Speed = 50
	Medium   Speed = 20
	Fast     Speed = 10
	VeryFast Speed = 5
	Instant  Speed = 0
)

type Terminal struct {
	width, height int
	padding       int
	frames        []frame
	current       *canvas.Canvas
	pos           image.Point
	face          font.Face
	theme         theme.Theme
	fontHeight    int
	showCursor    bool
	currentColour color.RGBA
}

type frame struct {
	img     image.Image
	delayMs int
}

func New(w, h int, face font.Face, theme theme.Theme) *Terminal {
	if face == nil {
		face = basicfont.Face7x13
	}
	t := &Terminal{
		width:         w,
		height:        h,
		face:          face,
		theme:         theme,
		padding:       15,
		fontHeight:    int(face.Metrics().Height >> 6),
		showCursor:    true,
		currentColour: theme.Foreground,
	}
	t.Clear()
	return t
}

func (t *Terminal) Clear() {
	t.current = canvas.New(t.width, t.height)
	t.current.Fill(t.theme.Background)
	t.pos = image.Point{
		X: t.padding,
		Y: t.padding,
	}
}

func (t *Terminal) SetHighlight(on bool) {
	if on {
		t.currentColour = t.theme.Highlight
	} else {
		t.currentColour = t.theme.Foreground
	}
}

func (t *Terminal) Frame(speed Speed) {
	if speed == Instant {
		return
	}
	clone := t.current.Clone()
	if t.showCursor {
		clone.DrawText(t.pos, t.theme.Foreground, t.face, "_")
	}
	fmt.Printf("Rendering frame %d...\n", len(t.frames))
	t.frames = append(t.frames, frame{
		img:     clone.Image(),
		delayMs: int(speed),
	})
}

func (t *Terminal) ShowCursor(show bool) {
	t.showCursor = show
}

func (t *Terminal) NewLine() {
	t.pos.X = t.padding
	newY := t.pos.Y + t.fontHeight
	// if we have room, jump down a line...
	if newY+t.fontHeight <= t.height-t.padding {
		t.pos.Y = newY
		return
	}

	// we need to scroll up...
	clone := t.current.Clone()
	t.current.Fill(t.theme.Background)
	t.current.DrawImage(image.Point{X: 0, Y: t.fontHeight}, clone.Image())
}

func (t *Terminal) Print(input string) {
	t.Type(input, Instant)
}

func (t *Terminal) Println(input string) {
	t.Print(input)
	t.NewLine()
}

func (t *Terminal) DrawImage(r image.Rectangle, img image.Image) {
	t.current.DrawImageAtRect(r, img)
}

func (t *Terminal) Type(input string, speed Speed) {
	for _, r := range input {
		if r == '\n' {
			t.NewLine()
			t.Frame(speed)
			continue
		}
		w, _ := t.current.DrawText(t.pos, t.currentColour, t.face, string(r))
		t.pos.X += w
		t.Frame(speed)
	}
}

func (t *Terminal) CursorToPx(x, y int) {
	t.pos = image.Point{X: x, Y: y}
}

func (t *Terminal) CursorToRow(row int) { // zero indexed
	t.pos = image.Point{X: t.padding, Y: t.padding + (row * t.fontHeight)}
}

func (t *Terminal) CursorToHome() { // zero indexed
	t.pos.X = t.padding
}

func (t *Terminal) CursorToLastRow() { // zero indexed
	t.CursorToRow(t.Rows() - 1)
}

func (t *Terminal) Rows() int {
	return (t.height - (t.padding * 2)) / t.fontHeight
}

func (t *Terminal) ClearLine() {
	t.current.Rect(0, t.pos.Y, t.width, t.pos.Y+t.fontHeight, t.theme.Background)
}

func (t *Terminal) ToGif(path string, loop bool) error {
	var frames []image.Image
	var delays []int
	fmt.Println("Preparing frames...")
	for _, f := range t.frames {
		frames = append(frames, f.img)
		delays = append(delays, f.delayMs)
	}
	fmt.Println("Creating gif...")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	var palettedFrames []*image.Paletted
	palette := color.Palette{
		t.theme.Background,
		t.theme.Foreground,
		t.theme.Highlight,
		color.RGBA{
			0x22, 0x22, 0x22, 0xff,
		},
		color.RGBA{
			0, 0, 0, 0xff,
		},
	}
	fmt.Println("Converting to paletted images...")
	for _, frame := range frames {
		//(&Canvas{img: frame.(*image.RGBA)}).ToPNG("test.png")
		paletted := image.NewPaletted(frame.Bounds(), palette)
		draw.Draw(paletted, frame.Bounds(), frame, image.Point{}, draw.Src)
		palettedFrames = append(palettedFrames, paletted)
	}
	fmt.Println("Encoding gif...")
	loopCount := -1
	if loop {
		loopCount = 0
	}
	return gif.EncodeAll(f, &gif.GIF{
		Image:     palettedFrames,
		Delay:     delays,
		LoopCount: loopCount,
	})
}
