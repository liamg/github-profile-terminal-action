package canvas

import (
    "image/png"
    "os"
)

func (c *Canvas) ToPNG(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() { _ = f.Close() }()
    return png.Encode(f, c.img)
}
