package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Text struct {
	font    *ttf.Font
	texture *sdl.Texture
	width   int32
	height  int32
}

func NewText(fontPath string, fontSize int) *Text {
	text := &Text{}
	font, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		panic(err)
		return nil
	}
	text.font = font;
	return text
}

func (text *Text) Width() int32 {
	return text.width
}

func (text *Text) Height() int32 {
	return text.height
}

func (text *Text) Update(content string, renderer *sdl.Renderer) {
	if text.font == nil {
		return
	}
	surface, err := text.font.RenderUTF8_Blended(content, sdl.Color{255, 255, 255, 255})
	if err == nil {
		text.width = surface.W
		text.height = surface.H
		if text.texture != nil {
			text.texture.Destroy()
			text.texture = nil
		}
		text.texture, err = renderer.CreateTextureFromSurface(surface)
		surface.Free()
		if err != nil {
			text.texture = nil
		}
	}
}

func (text *Text) Draw(renderer *sdl.Renderer, x int32, y int32) {
	if text.texture != nil {
		dst := &sdl.Rect{x, y, text.width, text.height}
		renderer.Copy(text.texture, nil, dst)
	}
}
