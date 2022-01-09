package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	ScreenWidth  int32 = 1920
	ScreenHeight int32 = 1080

	fontPath string = "resources/font/Share-TechMono.ttf"

	asteroidFontSize int = 20
	asteroidFont     *ttf.Font

	menuLogoFontSize int = 128
	menuItemFontSize int = 42
	menuItemSelected int = 0

	currentWordWidth    int32 = 350
	currentWordHeight   int32 = 37
	currentWordFontSize int   = 32
	currentWordMargin   int32 = 16
	currentWordPadding  int32 = 5
	currentWordBorder   int32 = 1
	currentWordFont     *ttf.Font

	currentWordTexture       *sdl.Texture
	currentWordTextureWidth  int32
	currentWordTextureHeight int32

	applicationRenderer *sdl.Renderer
	applicationRunning  bool
	gamePaused          bool
	gameOver            bool
	mainMenu            bool

	overlayGameOver *Text
	overlayScore    *Text
	overlayLevel    *Text
	hudEarth        *Text
	hudScore        *Text
	menuItemStart   *Text
	menuItemQuit    *Text
	menuLogo        *Text

	hudFontSize     int   = 32
	hudMarginRight  int32 = 16
	hudMarginBottom int32 = 8

	levelFontSize   int     = 92
	levelTimeToShow float32 = 2500.0
	levelTimeLeft   float32

	currentGame     *Game
	currentPlayer   *Player
	currentWord     string
	currentAsteroid *Asteroid

	playerScore int
)

func handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			applicationRunning = false
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYUP {
				continue
			}
			if t.Keysym.Sym == sdl.K_ESCAPE {
				if mainMenu {
					applicationRunning = false
				} else {
					if len(currentWord) > 0 {
						currentWord = ""
						if currentAsteroid != nil {
							currentAsteroid.Untarget()
							currentAsteroid = nil
						}
						updateCurrentWordTexture()
					} else {
						mainMenu = true
						gameOver = false
					}
				}
			} else if t.Keysym.Sym == sdl.K_BACKSPACE {
				if len(currentWord) > 0 {
					index := len(currentWord) - 1
					currentWord = currentWord[:index]
					updateCurrentWordTexture()
				}
			} else if t.Keysym.Sym == sdl.K_UP {
				if mainMenu {
					if menuItemSelected == 0 {
						menuItemSelected = 1
					} else {
						menuItemSelected = 0
					}
					createMainMenu()
				}
			} else if t.Keysym.Sym == sdl.K_DOWN {
				if mainMenu {
					if menuItemSelected == 0 {
						menuItemSelected = 1
					} else {
						menuItemSelected = 0
					}
					createMainMenu()
				}
			} else if t.Keysym.Sym == sdl.K_RETURN {
				if mainMenu {
					if menuItemSelected == 0 {
						startGame()
						mainMenu = false
						gameOver = false
					} else if menuItemSelected == 1 {
						applicationRunning = false
					}
				}
			} else {
				if gameOver || gamePaused {
					return
				}
				key := int(t.Keysym.Sym)
				if key >= 97 && key <= 122 {
					character := string(rune(key))
					if len(currentWord) == 0 {
						asteroid := currentGame.GetMatchingAsteroid(character)
						if asteroid != nil {
							currentAsteroid = asteroid
							currentAsteroid.Target()
							currentWord += character
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
								updateCurrentWordTexture()
								if len(currentWord) == wordLen {
									currentAsteroid.Destroy()
									playerScore += (len(currentAsteroid.word) * currentGame.level) * 10
									hudScore.Update(fmt.Sprintf("Score: %d", playerScore), applicationRenderer)
									currentAsteroid = nil
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
		overlayScore.Update(fmt.Sprintf("Your score: %d", playerScore), applicationRenderer)
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

	var windowFlags uint32 = sdl.WINDOW_FULLSCREEN_DESKTOP

	window, err := sdl.CreateWindow("Astrotyper", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, 0, 0, windowFlags)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	applicationRenderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer applicationRenderer.Destroy()

	ScreenWidth, ScreenHeight, err = applicationRenderer.GetOutputSize()
	if err != nil {
		panic(err)
	}

	background1 := NewBackground(100, 1, 1, 0.2)
	background2 := NewBackground(10, 1, 1, 0.3)

	createMainMenu()

	currentTime := sdl.GetTicks()
	lastTime := currentTime
	var deltaTime float32

	mainMenu = true

	applicationRunning = true
	for applicationRunning {
		currentTime = sdl.GetTicks()
		deltaTime = float32(currentTime - lastTime)
		lastTime = currentTime

		handleEvents()

		if !gameOver {
			background1.Update(deltaTime)
			background2.Update(deltaTime)
		}

		if !mainMenu {
			if !gamePaused && !gameOver {
				currentPlayer.Update(deltaTime)
				currentGame.Update(deltaTime)
			}
		}

		applicationRenderer.SetDrawColor(0, 0, 0, 255)
		applicationRenderer.Clear()

		background1.Draw(applicationRenderer)
		background2.Draw(applicationRenderer)

		if !mainMenu {
			currentPlayer.Draw(applicationRenderer)
			currentGame.Draw(applicationRenderer)

			drawLevel(deltaTime)
			drawGameOver()
			drawHUD()
			drawCurrentWord()
		} else {
			drawMainMenu()
		}

		applicationRenderer.Present()
	}

	//levelFont.Close()

	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func startGame() {
	if asteroidFont == nil {
		asteroidFont = openFont(fontPath, asteroidFontSize)
	}
	if currentWordFont == nil {
		currentWordFont = openFont(fontPath, currentWordFontSize)
	}
	updateCurrentWordTexture()

	if hudEarth == nil {
		hudEarth = NewText(fontPath, hudFontSize)
	}
	hudEarth.Update("Earth: 100%", applicationRenderer)
	if hudScore == nil {
		hudScore = NewText(fontPath, hudFontSize)
	}
	hudScore.Update("Score: 0", applicationRenderer)

	if overlayLevel == nil {
		overlayLevel = NewText(fontPath, levelFontSize)
	}
	handleNextLevel(1)

	if overlayGameOver == nil {
		overlayGameOver = NewText(fontPath, levelFontSize)
	}
	overlayGameOver.Update("GAME OVER", applicationRenderer)
	if overlayScore == nil {
		overlayScore = NewText(fontPath, levelFontSize)
	}

	if currentPlayer == nil {
		currentPlayer = NewPlayer(applicationRenderer)
	}
	currentPlayer.Reset()
	if currentGame == nil {
		currentGame = NewGame()
	}
	currentGame.Start(handleAsteroidNotDestroyed, handleNextLevel)

	gameOver = false
	gamePaused = false
	playerScore = 0
}

func openFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(fontPath, size)
	if err != nil {
		panic(err)
	}
	return font
}

func createMainMenu() {
	if menuLogo == nil {
		menuLogo = NewText(fontPath, menuLogoFontSize)
		menuLogo.Update("Astrotyper", applicationRenderer)
	}
	if menuItemStart == nil {
		menuItemStart = NewText(fontPath, menuItemFontSize)
	}
	menuItemStartText := "New Game"
	if menuItemSelected == 0 {
		menuItemStartText = "* New Game *"
	}
	menuItemStart.Update(menuItemStartText, applicationRenderer)
	if menuItemQuit == nil {
		menuItemQuit = NewText(fontPath, menuItemFontSize)
	}
	menuItemQuitText := "Quit"
	if menuItemSelected == 1 {
		menuItemQuitText = "* Quit *"
	}
	menuItemQuit.Update(menuItemQuitText, applicationRenderer)
}

func drawMainMenu() {
	menuLogo.Draw(applicationRenderer, (ScreenWidth/2)-(menuLogo.Width()/2), 128)
	menuItemStart.Draw(applicationRenderer,
		(ScreenWidth/2)-(menuItemStart.Width()/2),
		(ScreenHeight/2)-(menuItemStart.Height()))
	menuItemQuit.Draw(applicationRenderer,
		(ScreenWidth/2)-(menuItemQuit.Width()/2),
		(ScreenHeight/2)+(menuItemStart.Height()))
}

func drawLevel(deltaTime float32) {
	if levelTimeLeft > 0.0 {
		overlayLevel.Draw(applicationRenderer,
			(ScreenWidth/2)-(overlayLevel.Width()/2),
			(ScreenHeight/3)-(overlayLevel.Height()/2))
		levelTimeLeft -= deltaTime
	}
}

func drawGameOver() {
	if gameOver {
		overlayGameOver.Draw(applicationRenderer,
			(ScreenWidth/2)-(overlayGameOver.Width()/2),
			(ScreenHeight/3)-(overlayGameOver.Height()/2))
		overlayScore.Draw(applicationRenderer,
			(ScreenWidth/2)-(overlayScore.Width()/2),
			(ScreenHeight/3)-(overlayScore.Height()/2)+overlayGameOver.Height()+64)
	}
}

func drawHUD() {
	hudEarth.Draw(applicationRenderer,
		ScreenWidth-hudEarth.Width()-hudMarginRight,
		ScreenHeight-hudEarth.Height()-hudMarginBottom-hudScore.Height())
	hudScore.Draw(applicationRenderer,
		ScreenWidth-hudScore.Width()-hudMarginRight,
		ScreenHeight-hudScore.Height()-hudMarginBottom)
}

func updateFontTexture(text string, font *ttf.Font, texture **sdl.Texture, width *int32, height *int32, color sdl.Color) {
	if texture != nil {
		t := *texture
		t.Destroy()
		*texture = nil
		*width = 0
		*height = 0
	}
	surface, err := font.RenderUTF8Blended(text, color)
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
	updateFontTexture(currentWord+"_",
		currentWordFont,
		&currentWordTexture,
		&currentWordTextureWidth,
		&currentWordTextureHeight,
		sdl.Color{
			R: 238,
			G: 232,
			B: 213,
			A: 255})
}

func drawCurrentWord() {
	background := &sdl.Rect{}
	border := &sdl.Rect{}

	background.X = (ScreenWidth / 2) - (currentWordWidth / 2) - currentWordPadding
	background.Y = ScreenHeight - currentWordHeight - currentWordPadding - currentWordMargin
	background.W = currentWordWidth + (currentWordPadding * 2)
	background.H = currentWordHeight + (currentWordPadding * 2)

	border.X = background.X - currentWordBorder
	border.Y = background.Y - currentWordBorder
	border.W = background.W + (currentWordBorder * 2)
	border.H = background.H + (currentWordBorder * 2)

	applicationRenderer.SetDrawColor(38, 139, 210, 255)
	applicationRenderer.FillRect(border)
	applicationRenderer.SetDrawColor(0, 43, 54, 255)
	applicationRenderer.FillRect(background)

	if currentWordTexture != nil {
		text := &sdl.Rect{}
		text.X = background.X + currentWordPadding
		text.Y = background.Y + currentWordPadding
		text.W = currentWordTextureWidth
		text.H = currentWordTextureHeight
		applicationRenderer.Copy(currentWordTexture, nil, text)
	}
}
