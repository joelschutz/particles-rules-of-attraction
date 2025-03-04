package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/particle"
)

type SimulationEngine struct {
	particlesGroups                                                                      []*particle.ParticleGroup
	rules                                                                                [][]Rule
	wrapped                                                                              bool
	maxEffectDistance, terminalVelocity, conservationOfMomentum, particleRepulsionFactor float64
}

// New is a Board constructor
func NewSimulationEngine(
	maxEffectDistance, terminalVelocity, conservationOfMomentum, particleRepulsionFactor float64,
	wrapped bool,
	rules [][]Rule,
	particleGroups []*particle.ParticleGroup) *SimulationEngine {
	se := new(SimulationEngine)

	se.particlesGroups = make([]*particle.ParticleGroup, 0)
	se.maxEffectDistance = maxEffectDistance
	se.terminalVelocity = terminalVelocity
	se.particleRepulsionFactor = particleRepulsionFactor
	se.conservationOfMomentum = conservationOfMomentum
	se.wrapped = wrapped
	se.rules = rules
	se.particlesGroups = particleGroups

	return se
}

// Setup prepares SimulationEngine
func (se *SimulationEngine) Setup() {
	for pgIndex := range se.particlesGroups {
		se.applyRule(pgIndex)
	}
}

// Clear removes all board particles
func (se *SimulationEngine) Clear() error {
	se.particlesGroups = nil
	se.rules = nil
	return nil
}

// Reset places particles back on initial positions
func (se *SimulationEngine) Reset() {
	for _, pg := range se.particlesGroups {
		pg.ResetPosition()
	}
}

func (se *SimulationEngine) NextFrame() *[]*particle.ParticleGroup {
	se.updateSimulation()
	return &se.particlesGroups
}

func (se *SimulationEngine) updateSimulation() {
	var rulesWg sync.WaitGroup
	rulesWg.Add(len(se.particlesGroups))

	for pgIndex := range se.particlesGroups {
		go func(i int) {
			defer rulesWg.Done()
			se.applyRule(i)
		}(pgIndex)
	}

	rulesWg.Wait()
}

func (se *SimulationEngine) applyRule(pg1Index int) {
	for i1, p1 := range se.particlesGroups[pg1Index].Particles {
		fx, fy := 0.0, 0.0
		for pg2Index, pl := range se.particlesGroups {
			g := se.getAttractionForceBetweenParticles(pg1Index, pg2Index)
			for i2, p2 := range pl.Particles {
				if i1 == i2 && pg1Index == pg2Index {
					continue
				}

				dx := float64(p1.GetX() - p2.GetX())
				if se.wrapped {
					dx2 := float64(p1.GetX() + (1 - p2.GetX()))
					if math.Abs(dx) > math.Abs(dx2) {
						dx = dx2
					}
				}

				dy := float64(p1.GetY() - p2.GetY())
				if se.wrapped {
					dy2 := float64(p1.GetY() + (1 - p2.GetY()))
					if math.Abs(dy) > math.Abs(dy2) {
						dy = dy2
					}
				}

				if dx != 0 || dy != 0 {
					d := dx*dx + dy*dy
					if d < se.maxEffectDistance {
						F := g / math.Sqrt(d)
						fx += F * dx
						fy += F * dy
					}
				} else {
					fx += (rand.Float64()*2 - 1) * se.particleRepulsionFactor
					fy += (rand.Float64()*2 - 1) * se.particleRepulsionFactor
				}
			}
		}

		// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// 	cursorPosX, cursorPosY := ebiten.CursorPosition()

		// 	g := -64.0

		// 	dx := float64(p1.GetX() - cursorPosX)
		// 	dy := float64(p1.GetY() - cursorPosY)

		// 	if dx != 0 || dy != 0 {
		// 		d := dx*dx + dy*dy
		// 		F := g / math.Sqrt(d)
		// 		fx += F * dx
		// 		fy += F * dy
		// 	}
		// } else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		// 	cursorPosX, cursorPosY := ebiten.CursorPosition()

		// 	g := 64.0

		// 	dx := float64(p1.GetX() - cursorPosX)
		// 	dy := float64(p1.GetY() - cursorPosY)

		// 	if dx != 0 || dy != 0 {
		// 		d := dx*dx + dy*dy
		// 		F := g / math.Sqrt(d)
		// 		fx += F * dx
		// 		fy += F * dy
		// 	}
		// }

		factor := se.conservationOfMomentum

		p1.Vx = (p1.Vx + fx) * factor
		if math.Abs(p1.Vx) > float64(se.terminalVelocity) {
			negativeX := math.Signbit(p1.Vx)
			p1.Vx = se.terminalVelocity
			if negativeX {
				p1.Vx *= -1
			}
		}
		if p1.Vx != 0 {
			newP1X := p1.GetX() + p1.Vx
			if newP1X <= 0 {
				if se.wrapped {
					newP1X += 1
				} else {
					p1.Vx *= -1
					newP1X = 0
				}
			} else if newP1X >= 1 {
				if se.wrapped {
					newP1X -= 1
				} else {
					p1.Vx *= -1
					newP1X = 1
				}
			}
			p1.SetX(newP1X)
		}

		p1.Vy = (p1.Vy + fy) * factor
		if math.Abs(p1.Vy) > float64(se.terminalVelocity) {
			negativeY := math.Signbit(p1.Vy)
			p1.Vy = se.terminalVelocity
			if negativeY {
				p1.Vy *= -1
			}
		}
		if p1.Vy != 0 {
			newP1Y := p1.GetY() + p1.Vy
			if newP1Y <= 0 {
				if se.wrapped {
					newP1Y += 1
				} else {
					p1.Vy *= -1
					newP1Y = 0
				}
			} else if newP1Y >= 1 {
				if se.wrapped {
					newP1Y -= 1
				} else {
					p1.Vy *= -1
					newP1Y = 1
				}
			}
			p1.SetY(newP1Y)
		}
		fmt.Println("Vels", fx, fy, p1.Vx, p1.Vy)
	}
}

func (se *SimulationEngine) getAttractionForceBetweenParticles(pg1Index, pg2Index int) float64 {
	return float64(se.rules[pg1Index][pg2Index])
}
