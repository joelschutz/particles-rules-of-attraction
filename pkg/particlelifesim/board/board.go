package board

import (
	image "image/color"
	"math"
	"math/rand"
	"sync"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/color"
	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/particle"
	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/rule"
)

// Board encapsulates simulation logic
type Board struct {
	particlesGroups []*particle.ParticleGroup

	width  int
	height int

	rules [][]rule.Rule

	paused    bool
	forwarded bool
	// reversed  bool

	wrapped                             bool
	maxEffectDistance, terminalVelocity float64
}

// New is a Board constructor
func New(w, h int, med, tv float64, wr bool) *Board {
	b := new(Board)

	b.width = w
	b.height = h
	b.particlesGroups = make([]*particle.ParticleGroup, 0)
	b.maxEffectDistance = med
	b.terminalVelocity = tv
	b.wrapped = wr

	return b
}

func (b *Board) randomX() int {
	return rand.Intn(b.width-50) + 25
}

func (b *Board) randomY() int {
	return rand.Intn(b.height-50) + 25
}

func (b *Board) createParticles(name string, numberOfParticles int, color image.Color) {
	pg := particle.NewParticleGroup(name, color)

	for i := 0; i < numberOfParticles; i++ {
		p := particle.New(b.randomX(), b.randomY())
		pg.Particles = append(pg.Particles, p)
	}
	b.particlesGroups = append(b.particlesGroups, pg)
}

// Setup prepares board
func (b *Board) Setup(numberOfParticles int) {
	b.createParticles("red", numberOfParticles, color.RED)
	b.createParticles("green", numberOfParticles, color.GREEN)
	b.createParticles("blue", numberOfParticles, color.BLUE)
	b.createParticles("yellow", numberOfParticles, color.YELLOW)
	b.createParticles("white", numberOfParticles, color.WHITE)
	b.createParticles("teal", numberOfParticles, color.TEAL)

	b.rules = rule.GenerateRandomAsymmetricRules(b.particlesGroups)
	b.paused = false
	for pgIndex := range b.particlesGroups {
		b.applyRule(pgIndex)
	}
}

// TogglePause toggles board pause
func (b *Board) TogglePause() {
	b.paused = !b.paused
}

// Forward sets forward
func (b *Board) Forward(forward bool) {
	b.forwarded = forward
}

// Update performs board updates
func (b *Board) Update() error {
	return nil
}

// Clear removes all board particles
func (b *Board) Clear() error {
	b.particlesGroups = nil
	b.rules = nil
	return nil
}

// Size returns board size
func (b *Board) Size() (w, h int) {
	return b.width, b.height
}

// Draw draws board
func (b *Board) Draw(boardImage *ebiten.Image) {
	b.drawParticles(boardImage)
}

func (b *Board) drawParticles(boardImage *ebiten.Image) {
	if !b.paused || b.forwarded {
		boardImage.Clear()
		b.updateSimulation()
		for _, pg := range b.particlesGroups {
			for _, p := range pg.Particles {
				boardImage.Set(p.X, p.Y, pg.Color)
			}
		}
	}
}

func (b *Board) updateSimulation() {
	var rulesWg sync.WaitGroup
	rulesWg.Add(len(b.particlesGroups))

	for pgIndex := range b.particlesGroups {
		go func(i int) {
			defer rulesWg.Done()
			b.applyRule(i)
		}(pgIndex)
	}

	rulesWg.Wait()
}

func (b *Board) applyRule(pg1Index int) {
	for i1, p1 := range b.particlesGroups[pg1Index].Particles {
		fx, fy := 0.0, 0.0
		for pg2Index, pl := range b.particlesGroups {
			g := b.getAttractionForceBetweenParticles(pg1Index, pg2Index)
			for i2, p2 := range pl.Particles {
				if i1 == i2 && pg1Index == pg2Index {
					continue
				}

				dx := float64(p1.X - p2.X)
				if b.wrapped {
					dx2 := float64(p1.X + (b.width - p2.X))
					if math.Abs(dx) > math.Abs(dx2) {
						dx = dx2
					}
				}

				dy := float64(p1.Y - p2.Y)
				if b.wrapped {
					dy2 := float64(p1.Y + (b.height - p2.Y))
					if math.Abs(dy) > math.Abs(dy2) {
						dy = dy2
					}
				}

				if dx != 0 || dy != 0 {
					d := dx*dx + dy*dy
					if d < b.maxEffectDistance {
						F := g / math.Sqrt(d)
						fx += F * dx
						fy += F * dy
					}
				} else {
					fx += rand.Float64()*2 - 1
					fy += rand.Float64()*2 - 1
				}
			}
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			cursorPosX, cursorPosY := ebiten.CursorPosition()

			g := -64.0

			dx := float64(p1.X - cursorPosX)
			dy := float64(p1.Y - cursorPosY)

			if dx != 0 || dy != 0 {
				d := dx*dx + dy*dy
				F := g / math.Sqrt(d)
				fx += F * dx
				fy += F * dy
			}
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			cursorPosX, cursorPosY := ebiten.CursorPosition()

			g := 64.0

			dx := float64(p1.X - cursorPosX)
			dy := float64(p1.Y - cursorPosY)

			if dx != 0 || dy != 0 {
				d := dx*dx + dy*dy
				F := g / math.Sqrt(d)
				fx += F * dx
				fy += F * dy
			}
		}

		factor := 0.1

		p1.Vx = (p1.Vx + fx) * factor
		if math.Abs(p1.Vx) > float64(b.terminalVelocity) {
			p1.Vx = b.terminalVelocity
			if math.Signbit(p1.Vx) {
				p1.Vx *= -1
			}
		}
		if p1.Vx != 0 {
			p1.X += int(p1.Vx)
			if p1.X <= 0 {
				if b.wrapped {
					p1.X += b.width
				} else {
					p1.Vx *= -1
					p1.X = 0
				}
			} else if p1.X >= b.width {
				if b.wrapped {
					p1.X -= b.width
				} else {
					p1.Vx *= -1
					p1.X = b.width - 1
				}
			}
		}

		p1.Vy = (p1.Vy + fy) * factor
		if math.Abs(p1.Vy) > float64(b.terminalVelocity) {
			p1.Vy = b.terminalVelocity
			if math.Signbit(p1.Vy) {
				p1.Vy *= -1
			}
		}
		if p1.Vy != 0 {
			p1.Y += int(p1.Vy)
			if p1.Y <= 0 {
				if b.wrapped {
					p1.Y += b.height
				} else {
					p1.Vy *= -1
					p1.Y = 0
				}
			} else if p1.Y >= b.height {
				if b.wrapped {
					p1.Y -= b.height
				} else {
					p1.Vy *= -1
					p1.Y = b.height - 1
				}
			}
		}
	}
}

func (b *Board) getAttractionForceBetweenParticles(pg1Index, pg2Index int) float64 {
	return float64(b.rules[pg1Index][pg2Index])
}
