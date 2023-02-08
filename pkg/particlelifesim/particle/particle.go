package particle

type Particle struct {
	x  float32
	y  float32
	Vx float32
	Vy float32
}

func NewParticle(x, y float32) *Particle {
	p := new(Particle)
	p.SetX(x)
	p.SetY(y)

	return p
}

// The values are normalized before storage
func (p *Particle) SetX(x float32) {
	p.x = x
}

func (p *Particle) SetY(y float32) {
	p.y = y
}

func (p *Particle) GetX() float32 {
	return p.x
}

func (p *Particle) GetY() float32 {
	return p.y
}
