package main

import (
	"math/rand"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	startNumberOfAsteroids     int = 5
	startDelayBetweenAsteroids float32 = 5000.0
	startAsteroidVelocity      float32 = 0.1
	startAsteroidY             float32 = -64.0
	
	asteroidsToSpawnIncrement      int = 1
	delayBetweenAsteroidsIncrement float32 = -500.0
	asteroidVelocityIncrement      float32 = 0.01
	
	minDelayBetweenAsteroids float32 = 500.0

	asteroidMinDamage int = 5
	asteroidMaxDamage int = 10
)

type Asteroid struct {
	rectangle sdl.Rect
	alive     bool
	x         float32
	y         float32
	velocity  float32
	word      string
}

type AsteroidNotDestroyed func(int)
type NextLevel            func(int)

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

func NewAsteroid(x, y, velocity float32) *Asteroid {
	asteroid := &Asteroid{}
	asteroid.alive = true
	asteroid.x = x
	asteroid.y = y
	asteroid.velocity = velocity
	return asteroid
}

func (asteroid *Asteroid) Damage() int {
	damage := rand.Intn(asteroidMaxDamage - asteroidMinDamage)
	damage += asteroidMinDamage
	return damage
}

func (asteroid *Asteroid) IsAlive() bool {
	return asteroid.alive
}

func (asteroid *Asteroid) Update(deltaTime float32) {
	asteroid.y += (asteroid.velocity * deltaTime);
	if asteroid.y > float32(ScreenHeight) {
		asteroid.alive = false
	}
}

func (asteroid *Asteroid) Draw(renderer *sdl.Renderer) {
	if asteroid.x < -64.0 || asteroid.x > float32(ScreenWidth) ||
		asteroid.y < -64.0 || asteroid.y > float32(ScreenHeight) {
		return;
	}
	asteroid.rectangle.X = int32(asteroid.x)
	asteroid.rectangle.Y = int32(asteroid.y)
	asteroid.rectangle.W = 64
	asteroid.rectangle.H = 64
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&asteroid.rectangle)
}

func NewGame() *Game {
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
}

func (game *Game) spawnNextAsteroid() {
	x := float32(rand.Intn(int(ScreenWidth)))
	asteroid := NewAsteroid(x, startAsteroidY, game.asteroidVelocity)
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
			} else {
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

