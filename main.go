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
	gamePaused          bool
	gameOver            bool
	
	overlayLevel *Text
	hudEarth     *Text

	hudFontSize     int = 32
	hudMarginRight  int32 = 16
	hudMarginBottom int32 = 8

	levelFontSize   int = 92
	levelTimeToShow float32 = 2500.0
	levelTimeLeft   float32

	currentPlayer *Player
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

func handleAsteroidNotDestroyed(damage int) {
	currentPlayer.TakeDamage(damage)
	text := fmt.Sprintf("Earth: %d%%", currentPlayer.CurrentHealth())
	hudEarth.Update(text, applicationRenderer)

	if currentPlayer.CurrentHealth() == 0 {
		gameOver = true
		levelTimeLeft = 0.0
	}
}

func handleNextLevel(level int) {
	text := fmt.Sprintf("Level %d", level)
	overlayLevel.Update(text, applicationRenderer)
	levelTimeLeft = levelTimeToShow
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

	hudEarth = NewText(fontPath, hudFontSize)
	hudEarth.Update("Earth: 100%", applicationRenderer)

	overlayLevel = NewText(fontPath, levelFontSize)
	handleNextLevel(1)

	overlayGameOver := NewText(fontPath, levelFontSize)
	overlayGameOver.Update("GAME OVER", applicationRenderer)

	background1 := NewBackground(100, 1, 1, 0.2)
	background2 := NewBackground(10, 1, 1, 0.3)
	background3 := NewBackground(1, 4, 4, 0.4)
	currentPlayer = NewPlayer(applicationRenderer)
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

		if !gamePaused && !gameOver {
			background1.Update(deltaTime)
			background2.Update(deltaTime)
			background3.Update(deltaTime)
			currentPlayer.Update(deltaTime)
			game.Update(deltaTime)
		}

		applicationRenderer.SetDrawColor(0, 0, 0, 255)
		applicationRenderer.Clear()

		background1.Draw(applicationRenderer)
		background2.Draw(applicationRenderer)
		background3.Draw(applicationRenderer)
		currentPlayer.Draw(applicationRenderer)
		game.Draw(applicationRenderer)

		//hudEarth.Draw(applicationRenderer)

		if levelTimeLeft > 0.0 {
			overlayLevel.Draw(applicationRenderer,
				(ScreenWidth/2)-(overlayLevel.Width()/2),
				(ScreenHeight/3)-(overlayLevel.Height()/2))
			levelTimeLeft -= deltaTime
		}
		if gameOver {
			overlayGameOver.Draw(applicationRenderer,
				(ScreenWidth/2)-(overlayGameOver.Width()/2),
				(ScreenHeight/3)-(overlayGameOver.Height()/2))
		}
		hudEarth.Draw(applicationRenderer,
			ScreenWidth-hudEarth.Width()-hudMarginRight,
			ScreenHeight-hudEarth.Height()-hudMarginBottom)

		applicationRenderer.Present()
	}

	//levelFont.Close()

	ttf.Quit()
	img.Quit()
	sdl.Quit()
}
