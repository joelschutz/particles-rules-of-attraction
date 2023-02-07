package board

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/simulation"
)

// Board encapsulates simulation logic
type Board struct {
	width  int
	height int

	paused    bool
	forwarded bool
	// reversed  bool

	se *simulation.SimulationEngine
}

// New is a Board constructor
func New(w, h int, se *simulation.SimulationEngine) *Board {
	b := new(Board)

	b.width = w
	b.height = h
	b.se = se

	return b
}

// Setup prepares board
func (b *Board) Setup(numberOfParticles int) {
	b.se.Setup(numberOfParticles)
	b.paused = false
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
	return b.se.Clear()
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
		pgs := b.se.NextFrame()
		for _, pg := range *pgs {
			for _, p := range pg.Particles {
				boardImage.Set(int(p.GetX()*float64(b.height)), int(p.GetY()*float64(b.height)), pg.Color)
			}
		}
	}
}
