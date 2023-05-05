package gogl

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	ID uint32
}

func CreateTextureFromFile(file string, WrapMode, MagFilter, MinFilter int32, flip bool) (*Texture, error) {
	imageFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if flip {
		width := img.Bounds().Max.X
		height := img.Bounds().Max.Y
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				rgba.Set(j, i, img.At(j, img.Bounds().Max.Y-1-i))
			}
		}
	} else {
		draw.Draw(rgba, img.Bounds(), img, image.Pt(0, 0), draw.Src)
	}

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, WrapMode)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, WrapMode)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, MagFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, MinFilter)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Size().X), int32(img.Bounds().Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return &Texture{
		ID: tex,
	}, nil
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
