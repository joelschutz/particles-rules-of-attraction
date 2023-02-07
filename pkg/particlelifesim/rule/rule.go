package rule

import (
	"math/rand"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/particle"
)

type Rule float64

// GenerateRandomRules generates random rules for every pair of particles
func GenerateRandomRules(pgs []*particle.ParticleGroup) [][]Rule {
	rules := make([][]Rule, len(pgs))
	for i := range rules {
		rules[i] = make([]Rule, len(pgs))
	}

	for IndexG0 := range pgs {
		for IndexG1 := range pgs {
			rules[IndexG0][IndexG1] = Rule(rand.Float64()*2 - 1)
			rules[IndexG1][IndexG0] = rules[IndexG0][IndexG1]
		}
	}

	return rules
}
