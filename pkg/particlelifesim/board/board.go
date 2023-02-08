package board

import (
	"fmt"

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
func (b *Board) Setup() {
	b.se.Setup()
	b.paused = false
}

// Reset places particles back on initial positions
func (b *Board) Reset() {
	b.se.Reset()
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
func (b *Board) Draw(boardImage, rulesImage *ebiten.Image) {
	b.drawParticles(boardImage)
	b.drawRules(rulesImage)
}

func (b *Board) drawRules(rulesImage *ebiten.Image) {
	matrix := ebiten.NewImage(b.se.RuleSize())
	matrix.WritePixels(simulation.DrawRulesMatrix(*b.se.GetRules()))
	mw, mh := matrix.Size()
	img := ebiten.NewImage(mw+1, mh+1)
	for i, v := range *b.se.GetParticleGroups() {
		img.Set(i+1, 0, v.Color)
		img.Set(0, i+1, v.Color)
		fmt.Println("clr", v.Color)
	}
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(1, 1)
	img.DrawImage(matrix, op1)

	op := &ebiten.DrawImageOptions{}
	sw, sh := rulesImage.Size()
	bw, bh := img.Size()
	op.GeoM.Scale(float64(sw/bw), float64(sh/bh))
	rulesImage.DrawImage(img, op)
}

func (b *Board) drawParticles(boardImage *ebiten.Image) {
	if !b.paused || b.forwarded {
		boardImage.Clear()
		pgs := b.se.NextFrame()
		for _, pg := range *pgs {
			for _, p := range pg.Particles {
				boardImage.Set(int(p.GetX()*float32(b.height)), int(p.GetY()*float32(b.height)), pg.Color)
			}
		}
	}
}
