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

	asteroidFontSize int = 20
	asteroidFont     *ttf.Font

	currentWordFontSize int = 32
	currentWordMargin   int32 = 8
	currentWordPadding  int32 = 2
	currentWordBorder   int32 = 1
	currentWordFont     *ttf.Font

	currentWordTexture       *sdl.Texture
	currentWordTextureWidth  int32
	currentWordTextureHeight int32
	restWordTexture          *sdl.Texture
	restWordTextureWidth     int32
	restWordTextureHeight    int32

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

	currentGame     *Game
	currentPlayer   *Player
	currentWord     string
	currentAsteroid *Asteroid
)

func handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			applicationRunning = false
		case *sdl.KeyDownEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				if len(currentWord) > 0 {
					currentWord = ""
					currentAsteroid = nil
					fmt.Printf("current word: %s\n", currentWord)
					updateCurrentWordTexture()
				} else {
					applicationRunning = false
				}
			} else if t.Keysym.Sym == sdl.K_BACKSPACE {
				if len(currentWord) > 0 {
					index := len(currentWord) - 1
					currentWord = currentWord[:index]
					fmt.Printf("current word: %s\n", currentWord)
					updateCurrentWordTexture()
				}
			} else {
				key := int(t.Keysym.Sym)
				if key >= 97 && key <= 122 {
					character := string(key)
					if len(currentWord) == 0 {
						asteroid := currentGame.GetMatchingAsteroid(character)
						if asteroid != nil {
							currentAsteroid = asteroid
							currentWord += character
							fmt.Printf("current word: %s\n", currentWord)
							updateCurrentWordTexture()
						}
					} else {
						word := currentAsteroid.Word()
						wordLen := len(word)
						currentWordLen := len(currentWord)
						if currentWordLen < wordLen {
							nextValid := string(word[currentWordLen])
							if character == nextValid {
								currentWord += character
								fmt.Printf("current word: %s\n", currentWord)
								updateCurrentWordTexture()
								if len(currentWord) == wordLen {
									currentAsteroid.Destroy()
									currentWord = ""
									updateCurrentWordTexture()
								}
							}
						}
					}
				}
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

	asteroidFont = openFont(fontPath, asteroidFontSize)
	currentWordFont = openFont(fontPath, currentWordFontSize)

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
	currentGame = NewGame()

	currentGame.Start(handleAsteroidNotDestroyed, handleNextLevel)

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
			currentGame.Update(deltaTime)
		}

		applicationRenderer.SetDrawColor(0, 0, 0, 255)
		applicationRenderer.Clear()

		background1.Draw(applicationRenderer)
		background2.Draw(applicationRenderer)
		background3.Draw(applicationRenderer)
		currentPlayer.Draw(applicationRenderer)
		currentGame.Draw(applicationRenderer)

		drawLevel(deltaTime)
		drawGameOver(overlayGameOver)
		drawHUD()
		drawCurrentWord()

		applicationRenderer.Present()
	}

	//levelFont.Close()

	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func openFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(fontPath, size)
	if err != nil {
		panic(err)
	}
	return font
}

func drawLevel(deltaTime float32) {
	if levelTimeLeft > 0.0 {
		overlayLevel.Draw(applicationRenderer,
			(ScreenWidth/2)-(overlayLevel.Width()/2),
			(ScreenHeight/3)-(overlayLevel.Height()/2))
		levelTimeLeft -= deltaTime
	}
}

func drawGameOver(overlayGameOver *Text) {
	if gameOver {
		overlayGameOver.Draw(applicationRenderer,
			(ScreenWidth/2)-(overlayGameOver.Width()/2),
			(ScreenHeight/3)-(overlayGameOver.Height()/2))
	}
}

func drawHUD() {
	hudEarth.Draw(applicationRenderer,
		ScreenWidth-hudEarth.Width()-hudMarginRight,
		ScreenHeight-hudEarth.Height()-hudMarginBottom)
}

func updateFontTexture(text string, font *ttf.Font, texture **sdl.Texture, width *int32, height *int32, color sdl.Color) {
	if texture != nil {
		t := *texture
		t.Destroy()
		*texture = nil
		*width = 0
		*height = 0
	}
	surface, err := font.RenderUTF8_Blended(text, color)
	if err == nil {
		w := surface.W
		h := surface.H
		t, err := applicationRenderer.CreateTextureFromSurface(surface)
		surface.Free()
		if err == nil {
			*texture = t
			*width = w
			*height = h
		}
	}
}

func updateCurrentWordTexture() {
	var restWord string

	if currentAsteroid != nil {
		asteroidWord := currentAsteroid.Word()
		currentWordLen := len(currentWord)
		if currentWordLen < len(asteroidWord) {
			restWord = asteroidWord[currentWordLen:]
		}
	}

	updateFontTexture(currentWord,
		currentWordFont,
		&currentWordTexture,
		&currentWordTextureWidth,
		&currentWordTextureHeight,
		sdl.Color{255, 255, 0, 255})

	if len(restWord) > 0 {
		updateFontTexture(restWord,
			currentWordFont,
			&restWordTexture,
			&restWordTextureWidth,
			&restWordTextureHeight,
			sdl.Color{255, 255, 255, 255})
	} else {
		if restWordTexture != nil {
			restWordTexture.Destroy()
			restWordTexture = nil
			restWordTextureWidth = 0
			restWordTextureHeight = 0
		}
	}
}

func drawCurrentWord() {
	if currentWordTexture != nil {
		text := &sdl.Rect{}
		rest := &sdl.Rect{}
		background := &sdl.Rect{}
		border := &sdl.Rect{}

		text.X = (ScreenWidth / 2) - (currentWordTextureWidth / 2) - (restWordTextureWidth / 2)
		text.Y = ScreenHeight - currentWordTextureHeight - currentWordPadding - currentWordMargin - currentWordBorder
		text.W = currentWordTextureWidth
		text.H = currentWordTextureHeight

		rest.X = text.X + text.W
		rest.Y = text.Y
		rest.W = restWordTextureWidth
		rest.H = restWordTextureHeight

		background.X = text.X - currentWordPadding
		background.Y = text.Y - currentWordPadding
		background.W = text.W + rest.W + (currentWordPadding * 2)
		background.H = text.H + (currentWordPadding * 2)

		border.X = background.X - currentWordBorder
		border.Y = background.Y - currentWordBorder
		border.W = background.W + (currentWordBorder * 2)
		border.H = background.H + (currentWordBorder * 2)

		applicationRenderer.SetDrawColor(255, 0, 0, 255)
		applicationRenderer.FillRect(border)
		applicationRenderer.SetDrawColor(50, 50, 50, 255)
		applicationRenderer.FillRect(background)
		applicationRenderer.Copy(currentWordTexture, nil, text)
		if restWordTexture != nil {
			applicationRenderer.Copy(restWordTexture, nil, rest)
		}
	}
}