package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

var (
	playerTexturePath    string = "resources/player.png"
	playerTextureWidth   int32 = 64
	playerTextureHeight  int32 = 64
	playerOffsetY        int32 = -192
	playerJetBeamOffsetX int32 = 8
	playerJetBeamOffsetY int32 = 48
)

type Player struct {
	rectangle sdl.Rect
	texture   *sdl.Texture
	jetBeam   *JetBeamParticleEffect
}

func NewPlayer(renderer *sdl.Renderer) *Player {
	texture, err := img.LoadTexture(renderer, playerTexturePath)
	if err != nil {
		return nil
	}
	player := &Player{
		sdl.Rect{
			(ScreenWidth/2)-(playerTextureWidth/2),
			ScreenHeight+playerOffsetY,
			playerTextureWidth,
			playerTextureHeight,
		},
		texture,
		NewJetBeamParticleEffect(
			float32((ScreenWidth/2)-(playerTextureWidth/4)+playerJetBeamOffsetX),
			float32(ScreenHeight+playerOffsetY+playerJetBeamOffsetY)),
	}
	return player
}

func (player *Player) Draw(renderer *sdl.Renderer) {
	player.jetBeam.Draw(renderer)
	renderer.Copy(player.texture, nil, &player.rectangle)
}

func (player *Player) Update(deltaTime float32) {
	player.jetBeam.Update(deltaTime)
}

