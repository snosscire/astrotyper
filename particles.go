package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	jetBeamParticleEffectWidth  int = 10
	jetBeamParticleEffectHeight int = 16

	jetBeamParticleEffectParticleWidth     int32   = 8
	jetBeamParticleEffectParticleHeight    int32   = 8
	jetBeamParticleEffectParticleVelocityX float32 = 0.0
	jetBeamParticleEffectParticleVelocityY float32 = 0.2

	jetBeamParticleEffectYellowParticles            int     = 20
	jetBeamParticleEffectYellowParticleMinAliveTime float32 = 50
	jetBeamParticleEffectYellowParticleMaxAliveTime float32 = 100
	jetBeamParticleEffectYellowParticleRed          uint8   = 255
	jetBeamParticleEffectYellowParticleGreen        uint8   = 255
	jetBeamParticleEffectYellowParticleBlue         uint8   = 115

	jetBeamParticleEffectOrangeParticles            int     = 40
	jetBeamParticleEffectOrangeParticleMinAliveTime float32 = 100
	jetBeamParticleEffectOrangeParticleMaxAliveTime float32 = 150
	jetBeamParticleEffectOrangeParticleRed          uint8   = 255
	jetBeamParticleEffectOrangeParticleGreen        uint8   = 160
	jetBeamParticleEffectOrangeParticleBlue         uint8   = 40

	explosionParticleEffectWhiteParticles  int   = 500
	explosionParticleEffectYellowParticles int   = 1000
	explosionParticleEffectOrangeParticles int   = 2000
	explosionParticleEffectParticleWidth   int32 = 8
	explosionParticleEffectParticleHeight  int32 = 8

	explosionParticleEffectWhiteParticleAliveTime   float32 = 250.0
	explosionParticleEffectWhiteParticleRed         uint8   = 255
	explosionParticleEffectWhiteParticleGreen       uint8   = 255
	explosionParticleEffectWhiteParticleBlue        uint8   = 255
	explosionParticleEffectWhiteParticleMinVelocity int     = -8
	explosionParticleEffectWhiteParticleMaxVelocity int     = 8

	explosionParticleEffectYellowParticleAliveTime   float32 = 250.0
	explosionParticleEffectYellowParticleRed         uint8   = 255
	explosionParticleEffectYellowParticleGreen       uint8   = 255
	explosionParticleEffectYellowParticleBlue        uint8   = 155
	explosionParticleEffectYellowParticleMinVelocity int     = -12
	explosionParticleEffectYellowParticleMaxVelocity int     = 12

	explosionParticleEffectOrangeParticleAliveTime   float32 = 250.0
	explosionParticleEffectOrangeParticleRed         uint8   = 255
	explosionParticleEffectOrangeParticleGreen       uint8   = 160
	explosionParticleEffectOrangeParticleBlue        uint8   = 40
	explosionParticleEffectOrangeParticleMinVelocity int     = -16
	explosionParticleEffectOrangeParticleMaxVelocity int     = 16
)

type Particle struct {
	rectangle sdl.Rect
	alive     bool
	aliveTime float32
	timeLeft  float32
	red       uint8
	green     uint8
	blue      uint8
	originX   float32
	originY   float32
	currentX  float32
	currentY  float32
	velocityX float32
	velocityY float32
}

type JetBeamParticleEffect struct {
	originX         float32
	originY         float32
	yellowParticles []*Particle
	orangeParticles []*Particle
}

type ExplosionParticleEffect struct {
	whiteParticles  []*Particle
	yellowParticles []*Particle
	orangeParticles []*Particle
}

func (particle *Particle) IsAlive() bool {
	return particle.alive
}

func (particle *Particle) Reset() {
	particle.alive = true
	particle.timeLeft = particle.aliveTime
	particle.currentX = particle.originX
	particle.currentY = particle.originY
}

func (particle *Particle) Update(deltaTime float32) {
	if particle.alive == false {
		return
	}
	particle.timeLeft -= deltaTime
	if particle.timeLeft <= 0.0 {
		particle.alive = false
		return
	}
	particle.currentX += (particle.velocityX * deltaTime)
	particle.currentY += (particle.velocityY * deltaTime)
}

func (particle *Particle) Draw(renderer *sdl.Renderer) {
	particle.rectangle.X = int32(particle.currentX)
	particle.rectangle.Y = int32(particle.currentY)
	renderer.SetDrawColor(particle.red, particle.green, particle.blue, 255)
	renderer.FillRect(&particle.rectangle)
}

func NewJetBeamParticleEffect(x float32, y float32) *JetBeamParticleEffect {
	jetBeam := &JetBeamParticleEffect{}
	jetBeam.originX = x
	jetBeam.originY = y

	for i := 1; i <= jetBeamParticleEffectYellowParticles; i++ {
		particle := jetBeam.newYellowParticle(jetBeam.originX, jetBeam.originY)
		jetBeam.yellowParticles = append(jetBeam.yellowParticles, particle)
	}
	for i := 1; i <= jetBeamParticleEffectOrangeParticles; i++ {
		particle := jetBeam.newOrangeParticle(jetBeam.originX, jetBeam.originY)
		jetBeam.orangeParticles = append(jetBeam.orangeParticles, particle)
	}

	return jetBeam
}

func (jetBeam *JetBeamParticleEffect) newYellowParticle(x float32, y float32) *Particle {
	return jetBeam.newParticle(x, y,
		jetBeamParticleEffectYellowParticleMinAliveTime,
		jetBeamParticleEffectYellowParticleMaxAliveTime,
		jetBeamParticleEffectYellowParticleRed,
		jetBeamParticleEffectYellowParticleGreen,
		jetBeamParticleEffectYellowParticleBlue)
}

func (jetBeam *JetBeamParticleEffect) newOrangeParticle(x float32, y float32) *Particle {
	return jetBeam.newParticle(x, y,
		jetBeamParticleEffectOrangeParticleMinAliveTime,
		jetBeamParticleEffectOrangeParticleMaxAliveTime,
		jetBeamParticleEffectOrangeParticleRed,
		jetBeamParticleEffectOrangeParticleGreen,
		jetBeamParticleEffectOrangeParticleBlue)
}

func (jetBeam *JetBeamParticleEffect) newParticle(x float32, y float32,
	minAlive float32, maxAlive float32, red uint8, green uint8, blue uint8) *Particle {

	particleX := float32(rand.Intn(jetBeamParticleEffectWidth)) + x
	particleY := float32(rand.Intn(jetBeamParticleEffectHeight)) + y
	aliveTime := float32(rand.Intn(int(maxAlive-minAlive))) + minAlive

	particle := &Particle{}
	particle.rectangle.X = 0
	particle.rectangle.Y = 0
	particle.rectangle.W = jetBeamParticleEffectParticleWidth
	particle.rectangle.H = jetBeamParticleEffectParticleHeight
	particle.alive = true
	particle.aliveTime = aliveTime
	particle.timeLeft = aliveTime
	particle.red = red
	particle.green = green
	particle.blue = blue
	particle.originX = particleX
	particle.originY = particleY
	particle.currentX = particleX
	particle.currentY = particleY
	particle.velocityX = jetBeamParticleEffectParticleVelocityX
	particle.velocityY = jetBeamParticleEffectParticleVelocityY

	return particle
}

func (jetBeam *JetBeamParticleEffect) updateParticle(particle *Particle, deltaTime float32) {
	particle.Update(deltaTime)
	if particle.IsAlive() == false {
		particle.Reset()
	}
}

func (jetBeam *JetBeamParticleEffect) Update(deltaTime float32) {
	for _, particle := range jetBeam.yellowParticles {
		jetBeam.updateParticle(particle, deltaTime)
	}
	for _, particle := range jetBeam.orangeParticles {
		jetBeam.updateParticle(particle, deltaTime)
	}
}

func (jetBeam *JetBeamParticleEffect) Draw(renderer *sdl.Renderer) {
	for _, particle := range jetBeam.orangeParticles {
		particle.Draw(renderer)
	}
	for _, particle := range jetBeam.yellowParticles {
		particle.Draw(renderer)
	}
}

func NewExplosionParticleEffect(x, y float32) *ExplosionParticleEffect {
	effect := &ExplosionParticleEffect{}

	x += 32
	y += 32

	for i := 1; i <= explosionParticleEffectWhiteParticles; i++ {
		particle := effect.newWhiteParticle(x, y)
		effect.whiteParticles = append(effect.whiteParticles, particle)
	}
	for i := 1; i <= explosionParticleEffectYellowParticles; i++ {
		particle := effect.newYellowParticle(x, y)
		effect.yellowParticles = append(effect.yellowParticles, particle)
	}
	for i := 1; i <= explosionParticleEffectOrangeParticles; i++ {
		particle := effect.newOrangeParticle(x, y)
		effect.orangeParticles = append(effect.orangeParticles, particle)
	}

	return effect
}

func (explosion *ExplosionParticleEffect) randomVelocity(min int, max int) float32 {
	random := rand.Intn(max-min) + min
	velocity := float32(random) / 100.0
	return velocity
}

func (explosion *ExplosionParticleEffect) randomPosition(min int, max int) float32 {
	random := rand.Intn(max-min) + min
	position := float32(random)
	return position
}

func (explosion *ExplosionParticleEffect) newWhiteParticle(x, y float32) *Particle {
	particle := explosion.newParticle(x, y)
	particle.rectangle.W = explosionParticleEffectParticleWidth
	particle.rectangle.H = explosionParticleEffectParticleHeight
	particle.aliveTime = explosionParticleEffectWhiteParticleAliveTime
	particle.timeLeft = particle.aliveTime
	particle.red = explosionParticleEffectWhiteParticleRed
	particle.green = explosionParticleEffectWhiteParticleGreen
	particle.blue = explosionParticleEffectWhiteParticleBlue
	particle.velocityX = explosion.randomVelocity(explosionParticleEffectWhiteParticleMinVelocity, explosionParticleEffectWhiteParticleMaxVelocity)
	particle.velocityY = explosion.randomVelocity(explosionParticleEffectWhiteParticleMinVelocity, explosionParticleEffectWhiteParticleMaxVelocity)
	return particle
}

func (explosion *ExplosionParticleEffect) newYellowParticle(x, y float32) *Particle {
	particle := explosion.newParticle(x, y)
	particle.rectangle.W = explosionParticleEffectParticleWidth
	particle.rectangle.H = explosionParticleEffectParticleHeight
	particle.aliveTime = explosionParticleEffectYellowParticleAliveTime
	particle.timeLeft = particle.aliveTime
	particle.red = explosionParticleEffectYellowParticleRed
	particle.green = explosionParticleEffectYellowParticleGreen
	particle.blue = explosionParticleEffectYellowParticleBlue
	particle.velocityX = explosion.randomVelocity(explosionParticleEffectYellowParticleMinVelocity, explosionParticleEffectYellowParticleMaxVelocity)
	particle.velocityY = explosion.randomVelocity(explosionParticleEffectYellowParticleMinVelocity, explosionParticleEffectYellowParticleMaxVelocity)
	return particle
}

func (explosion *ExplosionParticleEffect) newOrangeParticle(x, y float32) *Particle {
	particle := explosion.newParticle(x, y)
	particle.rectangle.W = explosionParticleEffectParticleWidth
	particle.rectangle.H = explosionParticleEffectParticleHeight
	particle.aliveTime = explosionParticleEffectOrangeParticleAliveTime
	particle.timeLeft = particle.aliveTime
	particle.red = explosionParticleEffectOrangeParticleRed
	particle.green = explosionParticleEffectOrangeParticleGreen
	particle.blue = explosionParticleEffectOrangeParticleBlue
	particle.velocityX = explosion.randomVelocity(explosionParticleEffectOrangeParticleMinVelocity, explosionParticleEffectOrangeParticleMaxVelocity)
	particle.velocityY = explosion.randomVelocity(explosionParticleEffectOrangeParticleMinVelocity, explosionParticleEffectOrangeParticleMaxVelocity)
	return particle
}

func (explosion *ExplosionParticleEffect) newParticle(x, y float32) *Particle {
	x = explosion.randomPosition(int(x-16), int(x+16))
	y = explosion.randomPosition(int(y-16), int(y+16))
	particle := &Particle{}
	particle.rectangle.X = 0
	particle.rectangle.Y = 0
	particle.alive = true
	particle.originX = x
	particle.originY = y
	particle.currentX = x
	particle.currentY = y
	return particle
}

func (explosion *ExplosionParticleEffect) Update(deltaTime float32) {
	for _, particle := range explosion.whiteParticles {
		if particle.IsAlive() {
			particle.Update(deltaTime)
		}
	}
	for _, particle := range explosion.yellowParticles {
		if particle.IsAlive() {
			particle.Update(deltaTime)
		}
	}
	for _, particle := range explosion.orangeParticles {
		if particle.IsAlive() {
			particle.Update(deltaTime)
		}
	}
}

func (explosion *ExplosionParticleEffect) Draw(renderer *sdl.Renderer) {
	for _, particle := range explosion.orangeParticles {
		if particle.IsAlive() {
			particle.Draw(renderer)
		}
	}
	for _, particle := range explosion.yellowParticles {
		if particle.IsAlive() {
			particle.Draw(renderer)
		}
	}
	for _, particle := range explosion.whiteParticles {
		if particle.IsAlive() {
			particle.Draw(renderer)
		}
	}
}
