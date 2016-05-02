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
