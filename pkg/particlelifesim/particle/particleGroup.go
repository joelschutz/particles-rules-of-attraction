package particle

import (
	"image/color"
	"math"
	"math/rand"
)

type ParticleGroup struct {
	Name             string
	Color            color.Color
	Particles        []*Particle
	initialPositions []Particle
}

func NewParticleGroup(name string, numberOfParticles int, color color.Color, initialPositions []Particle) *ParticleGroup {
	pg := new(ParticleGroup)
	pg.Name = name
	pg.Color = color
	pg.Particles = placeParticles(numberOfParticles, initialPositions)
	pg.initialPositions = initialPositions

	return pg
}

func (pg *ParticleGroup) ResetPosition() {
	pg.Particles = placeParticles(len(pg.Particles), pg.initialPositions)
}

func placeParticles(n int, p []Particle) (ptcs []*Particle) {
	nClusters := len(p)
	if nClusters <= 0 {
		// Place particles randonly if no clusters are passed
		for i := 0; i < n; i++ {
			p := NewParticle(rand.Float64(), rand.Float64())
			ptcs = append(ptcs, p)
		}
	} else {
		// Place particles proportionally in each clusters
		for i := 0; i < n; i++ {
			tIndex := int(math.Mod(float64(i), float64(nClusters)))
			tParticle := p[tIndex]
			p := NewParticle(tParticle.GetX(), tParticle.GetY())
			ptcs = append(ptcs, p)
		}
	}
	return ptcs
}
