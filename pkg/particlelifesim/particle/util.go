package particle

var MIDDLE []Particle = []Particle{
	{.5, .5, 0, 0},
}

var CORNERS []Particle = []Particle{
	{0, 0, 0, 0},
	{1, 0, 0, 0},
	{0, 1, 0, 0},
	{1, 1, 0, 0},
}

var GRID_1 []Particle = []Particle{
	{.25, .25, 0, 0},
	{.75, .25, 0, 0},
	{.25, .75, 0, 0},
	{.75, .75, 0, 0},
}

var GRID_2 []Particle = []Particle{
	{.25, .25, 0, 0},
	{.25, .5, 0, 0},
	{.25, .75, 0, 0},
	{.75, .25, 0, 0},
	{.75, .5, 0, 0},
	{.75, .75, 0, 0},
}
