package particle

import (
	"image/color"
)

type ParticleGroup struct {
	Name      string
	Color     color.Color
	Particles []*Particle
}

func NewParticleGroup(name string, color color.Color) *ParticleGroup {
	pl := new(ParticleGroup)
	pl.Name = name
	pl.Color = color
	pl.Particles = make([]*Particle, 0)

	return pl
}
