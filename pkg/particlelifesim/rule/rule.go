package rule

import (
	"math/rand"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/particle"
)

type Rule float64

// GenerateRules creates 2D array to store particle rules
func GenerateRules(pgs []*particle.ParticleGroup) [][]Rule {
	rules := make([][]Rule, len(pgs))
	for i := range rules {
		rules[i] = make([]Rule, len(pgs))
	}
	return rules
}

// GenerateRandomSymmetricRules generates random rules for every pair of particles in a symmetric configuration
func GenerateRandomSymmetricRules(pgs []*particle.ParticleGroup) [][]Rule {
	rules := GenerateRules(pgs)

	for IndexG0 := range pgs {
		for IndexG1 := range pgs {
			rules[IndexG0][IndexG1] = Rule(rand.Float64()*2 - 1)
			rules[IndexG1][IndexG0] = rules[IndexG0][IndexG1]
		}
	}

	return rules
}

// GenerateRandomAsymmetricRules generates random rules for every pair of particles in a asymmetric configuration
func GenerateRandomAsymmetricRules(pgs []*particle.ParticleGroup) [][]Rule {
	rules := GenerateRules(pgs)

	for IndexG0 := range pgs {
		for IndexG1 := range pgs {
			rules[IndexG0][IndexG1] = Rule(rand.Float64()*2 - 1)
			rules[IndexG1][IndexG0] = Rule(rand.Float64()*2 - 1)
		}
	}

	return rules
}
