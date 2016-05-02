package main

import (
	"math/rand"
	"github.com/veandco/go-sdl2/sdl"
)

type Star struct {
	rectangle sdl.Rect
	x         float32
	y         float32
	velocity  float32
}

type Background struct {
	stars []*Star
}

func NewStar(width int32, height int32, velocity float32) *Star {
	x := float32(rand.Intn(int(ScreenWidth)))
	y := float32(rand.Intn(int(ScreenHeight)))
	star := &Star{
		sdl.Rect{
			0,
			0,
			width,
			height,
		},
		x,
		y,
		velocity,
	}
	return star
}

func (star *Star) Update(deltaTime float32) {
	star.y += (star.velocity * deltaTime);
	if star.y > float32(ScreenHeight) {
		star.y = -1.0
		star.x = float32(rand.Intn(int(ScreenWidth)))
	}
}

func (star *Star) Draw(renderer *sdl.Renderer) {
	if star.x < 0.0 || star.x > float32(ScreenWidth) ||
		star.y < 0.0 || star.y > float32(ScreenHeight) {
		return;
	}
	star.rectangle.X = int32(star.x)
	star.rectangle.Y = int32(star.y)
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(&star.rectangle)
}

func NewBackground(numberOfStars uint,
	starWidth int32,
	starHeight int32,
	starVelocity float32) *Background {

	background := &Background{}

	var i uint
	i = 1
	for ; i <= numberOfStars; i++ {
		star := NewStar(starWidth, starHeight, starVelocity)
		background.stars = append(background.stars, star)
	}
	
	return background
}


func (background *Background) Update(deltaTime float32) {
	for _, star := range background.stars {
		star.Update(deltaTime)
	}
}

func (background *Background) Draw(renderer *sdl.Renderer) {
	for _, star := range background.stars {
		star.Draw(renderer)
	}
}

