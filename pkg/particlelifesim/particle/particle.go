package particle

type Particle struct {
	x  float64
	y  float64
	Vx float64
	Vy float64
}

func NewParticle(x, y float64) *Particle {
	p := new(Particle)
	p.SetX(x)
	p.SetY(y)

	return p
}

// The values are normalized before storage
func (p *Particle) SetX(x float64) {
	p.x = x
}

func (p *Particle) SetY(y float64) {
	p.y = y
}

func (p *Particle) GetX() float64 {
	return p.x
}

func (p *Particle) GetY() float64 {
	return p.y
}
