package simulation

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/particle"
)

type Rule float32

func (r Rule) String() string {
	return strconv.FormatFloat(float64(r), 'f', 'f', 32)
}

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
		for IndexG1 := IndexG0; IndexG1 < len(pgs); IndexG1++ {
			rules[IndexG0][IndexG1] = Rule(rand.Float32()*2 - 1)
			if IndexG1 != IndexG0 {
				rules[IndexG1][IndexG0] = rules[IndexG0][IndexG1]
			}
		}
	}

	return rules
}

// GenerateRandomAsymmetricRules generates random rules for every pair of particles in a asymmetric configuration
func GenerateRandomAsymmetricRules(pgs []*particle.ParticleGroup) [][]Rule {
	rules := GenerateRules(pgs)

	for IndexG0 := range pgs {
		for IndexG1 := IndexG0; IndexG1 < len(pgs); IndexG1++ {
			rules[IndexG0][IndexG1] = Rule(rand.Float32()*2 - 1)
			if IndexG1 != IndexG0 {
				rules[IndexG1][IndexG0] = Rule(rand.Float32()*2 - 1)
			}
		}
	}

	return rules
}

func DrawRulesMatrix(r [][]Rule) []byte {
	arr := []byte{}
	for IndexG0 := range r {
		for IndexG1 := range r[IndexG0] {
			c := []byte{0, 0, 0, 255}
			v := float32(r[IndexG0][IndexG1])
			if math.Signbit(float64(v)) {
				c[0] = byte(math.Abs(float64(v)) * 255)
			} else {
				c[2] = byte(math.Abs(float64(v)) * 255)
			}
			arr = append(arr, c...)
		}
	}
	return arr
}

func ExportRulesAsJson(r [][]Rule) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}

func SaveRulesAsCsv(r [][]Rule, filename string) error {
	f, err := os.Create(filename + ".csv")
	defer f.Close()

	if err != nil {
		log.Println("error writing record to file", err)
		return err
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range r {
		line := []string{}
		for _, v := range record {
			line = append(line, v.String())
		}
		if err := w.Write(line); err != nil {
			log.Println("error writing record to file", err)
			return err
		}
	}
	return nil
}
