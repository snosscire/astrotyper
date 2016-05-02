package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	ScreenWidth  int32 = 1920
	ScreenHeight int32 = 1080

	fontPath string = "resources/font/Share-TechMono.ttf"
	
	applicationRenderer *sdl.Renderer
	applicationRunning  bool
	
	levelFontSize      int = 92
	levelFont          *ttf.Font
	levelTexture       *sdl.Texture
	levelTextureWidth  int32
	levelTextureHeight int32
	levelTimeToShow    float32 = 2500.0
	levelTimeLeft      float32
)

func handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			applicationRunning = false
		case *sdl.KeyDownEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				applicationRunning = false
			}
		}
	}
}

func handleAsteroidNotDestroyed() {
	fmt.Println("Asteroid not destroyed!")
}

func handleNextLevel(level int) {
	text := fmt.Sprintf("Level %d", level)
	surface, err := levelFont.RenderUTF8_Blended(text, sdl.Color{255, 255, 255, 255})
	if err == nil {
		levelTextureWidth = surface.W
		levelTextureHeight = surface.H
		if levelTexture != nil {
			levelTexture.Destroy()
			levelTexture = nil
		}
		levelTexture, err = applicationRenderer.CreateTextureFromSurface(surface)
		surface.Free()
		if levelTexture != nil && err == nil {
			levelTimeLeft = levelTimeToShow
		}
	}
}

func init() {
	runtime.LockOSThread()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	sdl.Init(sdl.INIT_EVERYTHING)
	img.Init(img.INIT_PNG)
	ttf.Init()

	var windowFlags uint32 = sdl.WINDOW_SHOWN

	if runtime.GOOS == `darwin` {
		ScreenWidth = 1280
		ScreenHeight = 800
		windowFlags = sdl.WINDOW_FULLSCREEN
	}

	window, err := sdl.CreateWindow("Astrotyper", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int(ScreenWidth), int(ScreenHeight), windowFlags)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	applicationRenderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer applicationRenderer.Destroy()

	levelFont, err = ttf.OpenFont(fontPath, levelFontSize)
	if err != nil {
		panic(err)
	}
	handleNextLevel(1)

	background1 := NewBackground(100, 1, 1, 0.2)
	background2 := NewBackground(10, 1, 1, 0.3)
	background3 := NewBackground(1, 4, 4, 0.4)
	player := NewPlayer(applicationRenderer)
	game := NewGame()

	game.Start(handleAsteroidNotDestroyed, handleNextLevel)

	currentTime := sdl.GetTicks()
	lastTime := currentTime
	var deltaTime float32

	applicationRunning = true
	for applicationRunning {
		currentTime = sdl.GetTicks()
		deltaTime = float32(currentTime - lastTime)
		lastTime = currentTime

		handleEvents()

		background1.Update(deltaTime)
		background2.Update(deltaTime)
		background3.Update(deltaTime)
		player.Update(deltaTime)
		game.Update(deltaTime)

		applicationRenderer.SetDrawColor(0, 0, 0, 255)
		applicationRenderer.Clear()

		background1.Draw(applicationRenderer)
		background2.Draw(applicationRenderer)
		background3.Draw(applicationRenderer)
		player.Draw(applicationRenderer)
		game.Draw(applicationRenderer)

		if levelTexture != nil && levelTimeLeft > 0.0 {
			rect := &sdl.Rect{
				(ScreenWidth/2)-(levelTextureWidth/2),
				(ScreenHeight/3)-(levelTextureHeight/2),
				levelTextureWidth,
				levelTextureHeight}
			applicationRenderer.Copy(levelTexture, nil, rect)
			levelTimeLeft -= deltaTime
		}

		applicationRenderer.Present()
	}

	levelFont.Close()

	ttf.Quit()
	img.Quit()
	sdl.Quit()
}
