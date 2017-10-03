package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	playerStartHealth    int    = 100
	playerTexturePath    string = "resources/player.png"
	playerTextureWidth   int32  = 64
	playerTextureHeight  int32  = 64
	playerOffsetY        int32  = -192
	playerJetBeamOffsetX int32  = 8
	playerJetBeamOffsetY int32  = 48
)

type Player struct {
	health    int
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
		playerStartHealth,
		sdl.Rect{
			X: (ScreenWidth / 2) - (playerTextureWidth / 2),
			Y: ScreenHeight + playerOffsetY,
			W: playerTextureWidth,
			H: playerTextureHeight,
		},
		texture,
		NewJetBeamParticleEffect(
			float32((ScreenWidth/2)-(playerTextureWidth/4)+playerJetBeamOffsetX),
			float32(ScreenHeight+playerOffsetY+playerJetBeamOffsetY)),
	}
	return player
}

func (player *Player) Reset() {
	player.health = playerStartHealth
}

func (player *Player) TakeDamage(damage int) {
	player.health -= damage
	if player.health < 0 {
		player.health = 0
	}
}

func (player *Player) CurrentHealth() int {
	return player.health
}

func (player *Player) Draw(renderer *sdl.Renderer) {
	player.jetBeam.Draw(renderer)
	renderer.Copy(player.texture, nil, &player.rectangle)
}

func (player *Player) Update(deltaTime float32) {
	player.jetBeam.Update(deltaTime)
}
