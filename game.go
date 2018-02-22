package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	startNumberOfAsteroids     int     = 5
	startDelayBetweenAsteroids float32 = 2000.0
	startAsteroidVelocity      float32 = 0.1
	startAsteroidY             float32 = -64.0

	asteroidsToSpawnIncrement      int     = 1
	delayBetweenAsteroidsIncrement float32 = -100.0
	asteroidVelocityIncrement      float32 = 0.01

	asteroidRegularWordColor  sdl.Color = sdl.Color{R: 220, G: 50, B: 47, A: 255}
	asteroidTargetedWordColor sdl.Color = sdl.Color{R: 133, G: 153, B: 0, A: 255}

	minDelayBetweenAsteroids float32 = 1000.0

	asteroidMinDamage   int   = 5
	asteroidMaxDamage   int   = 10
	asteroidWordMargin  int32 = 10
	asteroidWordPadding int32 = 1
	asteroidWordBorder  int32 = 1

	asteroid1TexturePath string = "resources/asteroid1.png"
	asteroid2TexturePath string = "resources/asteroid2.png"
	asteroid3TexturePath string = "resources/asteroid3.png"
	asteroid4TexturePath string = "resources/asteroid4.png"
	asteroidTextures     []*sdl.Texture

	wordList []string
)

type Asteroid struct {
	rectangle         sdl.Rect
	alive             bool
	destroyed         bool
	targeted          bool
	x                 float32
	y                 float32
	velocity          float32
	texture           *sdl.Texture
	word              string
	wordTexture       *sdl.Texture
	wordTextureWidth  int32
	wordTextureHeight int32
	explosion         *ExplosionParticleEffect
}

type AsteroidNotDestroyed func(int)
type NextLevel func(int)

type Game struct {
	asteroids                  []*Asteroid
	level                      int
	numberOfAsteroidsToSpawn   int
	asteroidsLeftToSpawn       int
	delayBetweenAsteroids      float32
	timeUntilNextAsteroidSpawn float32
	asteroidVelocity           float32
	asteroidNotDestroyed       AsteroidNotDestroyed
	nextLevel                  NextLevel
}

func NewAsteroid(x, y, velocity float32, level int) *Asteroid {
	asteroid := &Asteroid{}
	asteroid.alive = true
	asteroid.destroyed = false
	asteroid.x = x
	asteroid.y = y
	asteroid.velocity = velocity
	asteroid.targeted = false
	asteroid.texture = randomAsteroidTexture()
	asteroid.word = randomWord(level)
	asteroid.updateWordTexture()
	return asteroid
}

func (asteroid *Asteroid) Word() string {
	return asteroid.word
}

func (asteroid *Asteroid) Damage() int {
	damage := rand.Intn(asteroidMaxDamage - asteroidMinDamage)
	damage += asteroidMinDamage
	return damage

}
func (asteroid *Asteroid) Destroy() {
	asteroid.alive = false
	asteroid.destroyed = true
	asteroid.explosion = NewExplosionParticleEffect(asteroid.x, asteroid.y)
}

func (asteroid *Asteroid) updateWordTexture() {
	color := asteroidRegularWordColor
	if asteroid.targeted {
		color = asteroidTargetedWordColor
	}
	surface, err := asteroidFont.RenderUTF8Blended(asteroid.word, color)
	if err == nil {
		asteroid.wordTextureWidth = surface.W
		asteroid.wordTextureHeight = surface.H
		asteroid.wordTexture, err = applicationRenderer.CreateTextureFromSurface(surface)
		surface.Free()
		if err != nil {
			asteroid.wordTexture = nil
		}
	}
}

func (asteroid *Asteroid) Target() {
	asteroid.targeted = true
	asteroid.updateWordTexture()
}

func (asteroid *Asteroid) Untarget() {
	asteroid.targeted = false
	asteroid.updateWordTexture()
}

func (asteroid *Asteroid) IsAlive() bool {
	if !asteroid.alive && asteroid.explosion != nil {
		return asteroid.explosion.IsAlive()
	}
	return asteroid.alive
}

func (asteroid *Asteroid) WasDestroyed() bool {
	return asteroid.destroyed
}

func (asteroid *Asteroid) Update(deltaTime float32) {
	if asteroid.alive {
		asteroid.y += (asteroid.velocity * deltaTime)
		if asteroid.y > float32(ScreenHeight) {
			asteroid.alive = false
		}
	} else {
		if asteroid.explosion != nil {
			asteroid.explosion.Update(deltaTime)
		}
	}
}

func (asteroid *Asteroid) Draw(renderer *sdl.Renderer) {
	if !asteroid.alive {
		if asteroid.explosion != nil {
			asteroid.explosion.Draw(renderer)
			return
		}
	}
	if asteroid.x < -64.0 || asteroid.x > float32(ScreenWidth) ||
		asteroid.y < -64.0 || asteroid.y > float32(ScreenHeight) {
		return
	}
	asteroid.rectangle.X = int32(asteroid.x)
	asteroid.rectangle.Y = int32(asteroid.y)
	asteroid.rectangle.W = 64
	asteroid.rectangle.H = 64
	renderer.Copy(asteroid.texture, nil, &asteroid.rectangle)

	if asteroid.wordTexture != nil {
		var wordX, wordY int32
		var bgX, bgY, bgW, bgH int32
		var borderX, borderY, borderW, borderH int32
		wordX = asteroid.rectangle.X + asteroid.rectangle.W + asteroidWordMargin
		wordY = asteroid.rectangle.Y + (asteroid.rectangle.H / 2) - (asteroid.wordTextureHeight / 2)
		bgX = wordX - asteroidWordPadding
		bgY = wordY - asteroidWordPadding
		bgW = asteroid.wordTextureWidth + (asteroidWordPadding * 2)
		bgH = asteroid.wordTextureHeight + (asteroidWordPadding * 2)
		borderX = bgX - asteroidWordBorder
		borderY = bgY - asteroidWordBorder
		borderW = bgW + (asteroidWordBorder * 2)
		borderH = bgH + (asteroidWordBorder * 2)

		borderColor := asteroidRegularWordColor
		if asteroid.targeted {
			borderColor = asteroidTargetedWordColor
		}
		renderer.SetDrawColor(borderColor.R, borderColor.G, borderColor.B, 255)
		renderer.FillRect(&sdl.Rect{
			X: borderX,
			Y: borderY,
			W: borderW,
			H: borderH,
		})
		renderer.SetDrawColor(0, 43, 54, 255)
		renderer.FillRect(&sdl.Rect{
			X: bgX,
			Y: bgY,
			W: bgW,
			H: bgH,
		})
		renderer.Copy(
			asteroid.wordTexture,
			nil,
			&sdl.Rect{
				X: wordX,
				Y: wordY,
				W: asteroid.wordTextureWidth,
				H: asteroid.wordTextureHeight,
			},
		)
	}
}

func NewGame() *Game {
	err := loadAsteroidTextures()
	if err != nil {
		panic(err)
	}

	game := &Game{}
	return game
}

func (game *Game) Start(asteroidNotDestroyed AsteroidNotDestroyed, nextLevel NextLevel) {
	game.level = 1
	game.numberOfAsteroidsToSpawn = startNumberOfAsteroids
	game.asteroidsLeftToSpawn = game.numberOfAsteroidsToSpawn
	game.delayBetweenAsteroids = startDelayBetweenAsteroids
	game.timeUntilNextAsteroidSpawn = game.delayBetweenAsteroids
	game.asteroidVelocity = startAsteroidVelocity
	game.asteroidNotDestroyed = asteroidNotDestroyed
	game.nextLevel = nextLevel
	game.asteroids = make([]*Asteroid, 0)
}

func (game *Game) GetMatchingAsteroid(firstCharacter string) *Asteroid {
	if len(game.asteroids) > 0 {
		for _, asteroid := range game.asteroids {
			if asteroid.IsAlive() {
				if firstCharacter == string(asteroid.word[0]) {
					return asteroid
				}
			}
		}
	}
	return nil
}

func (game *Game) spawnNextAsteroid() {
	x := float32(rand.Intn(int(ScreenWidth)-512) + 64)
	asteroid := NewAsteroid(x, startAsteroidY, game.asteroidVelocity, game.level)
	game.asteroids = append(game.asteroids, asteroid)
	game.asteroidsLeftToSpawn--
}

func (game *Game) goToNextLevel() {
	game.asteroids = make([]*Asteroid, 0)
	game.level++
	game.numberOfAsteroidsToSpawn += asteroidsToSpawnIncrement
	game.asteroidsLeftToSpawn = game.numberOfAsteroidsToSpawn
	game.delayBetweenAsteroids += delayBetweenAsteroidsIncrement
	if game.delayBetweenAsteroids < minDelayBetweenAsteroids {
		game.delayBetweenAsteroids = minDelayBetweenAsteroids
	}
	game.timeUntilNextAsteroidSpawn = game.delayBetweenAsteroids
	game.asteroidVelocity += asteroidVelocityIncrement
	if game.nextLevel != nil {
		game.nextLevel(game.level)
	}
}

func (game *Game) Update(deltaTime float32) {
	if game.asteroidsLeftToSpawn > 0 {
		game.timeUntilNextAsteroidSpawn -= deltaTime
		if game.timeUntilNextAsteroidSpawn <= 0.0 {
			game.spawnNextAsteroid()

			var leftOverTime float32 = 0.0
			if game.timeUntilNextAsteroidSpawn < 0.0 {
				leftOverTime = game.timeUntilNextAsteroidSpawn
			}
			game.timeUntilNextAsteroidSpawn = game.delayBetweenAsteroids + leftOverTime
		}
	}

	allAsteroidsDead := true
	for _, asteroid := range game.asteroids {
		if asteroid.IsAlive() {
			asteroid.Update(deltaTime)
			if asteroid.IsAlive() {
				allAsteroidsDead = false
			} else if !asteroid.WasDestroyed() {
				if game.asteroidNotDestroyed != nil {
					game.asteroidNotDestroyed(asteroid.Damage())
				}
			}
		}
	}

	if allAsteroidsDead && game.asteroidsLeftToSpawn <= 0 {
		game.goToNextLevel()
	}
}

func (game *Game) Draw(renderer *sdl.Renderer) {
	for _, asteroid := range game.asteroids {
		if asteroid.IsAlive() {
			asteroid.Draw(renderer)
		}
	}
}

func createWordList() {
	wordList = append(wordList, "car")
	wordList = append(wordList, "eat")
	wordList = append(wordList, "fat")
	wordList = append(wordList, "gun")
	wordList = append(wordList, "hug")
	wordList = append(wordList, "net")
	wordList = append(wordList, "put")
	wordList = append(wordList, "war")

	wordList = append(wordList, "five")
	wordList = append(wordList, "four")
	wordList = append(wordList, "nine")
	wordList = append(wordList, "bear")
	wordList = append(wordList, "food")
	wordList = append(wordList, "last")
	wordList = append(wordList, "fast")
	wordList = append(wordList, "port")
	wordList = append(wordList, "door")

	wordList = append(wordList, "seven")
	wordList = append(wordList, "eight")
	wordList = append(wordList, "right")
	wordList = append(wordList, "smite")
	wordList = append(wordList, "queue")
	wordList = append(wordList, "smart")
	wordList = append(wordList, "smear")
	wordList = append(wordList, "dance")
	wordList = append(wordList, "blast")

	wordList = append(wordList, "eleven")
	wordList = append(wordList, "twelve")
	wordList = append(wordList, "tought")
	wordList = append(wordList, "bought")
	wordList = append(wordList, "trench")
	wordList = append(wordList, "cought")
	wordList = append(wordList, "faster")
	wordList = append(wordList, "answer")
	wordList = append(wordList, "slower")

	wordList = append(wordList, "monster")
	wordList = append(wordList, "bouncer")
	wordList = append(wordList, "assault")
	wordList = append(wordList, "message")
	wordList = append(wordList, "corrupt")
	wordList = append(wordList, "acquire")
	wordList = append(wordList, "explodes")
	wordList = append(wordList, "contains")
	wordList = append(wordList, "tailoring")

	wordList = append(wordList, "sacrifice")
	wordList = append(wordList, "feedback")
	wordList = append(wordList, "purchase")
	wordList = append(wordList, "financial")
	wordList = append(wordList, "difficult")
	wordList = append(wordList, "department")
	wordList = append(wordList, "exchange")
	wordList = append(wordList, "exhibiting")
	wordList = append(wordList, "dedication")
	wordList = append(wordList, "complicated")
}

func init() {
	createWordList()
}

func loadAsteroidTextures() error {
	texturePaths := []string{
		asteroid1TexturePath,
		asteroid2TexturePath,
		asteroid3TexturePath,
		asteroid4TexturePath,
	}
	for _, texturePath := range texturePaths {
		texture, err := img.LoadTexture(applicationRenderer, texturePath)
		if err != nil {
			return err
		}
		asteroidTextures = append(asteroidTextures, texture)
	}
	return nil
}

func randomAsteroidTexture() *sdl.Texture {
	max := len(asteroidTextures) - 1
	random := rand.Intn(max)
	return asteroidTextures[random]
}

func randomWord(level int) string {
	level--
	if level > 20 {
		level = 20
	}
	max := len(wordList) - 1
	random := rand.Intn(max-level) + level
	return wordList[random]
}
