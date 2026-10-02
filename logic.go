package main

import (
	"math/rand"
)

func CanICatchIt(baseExperience int) bool {
	chance := rand.Float64()

	if baseExperience <= 10 {
		chance *= 1.25
	} else if baseExperience <= 100 {
		chance *= 1.1
	} else if baseExperience <= 1000 {
		chance *= 1.01
	}

	return chance > 0.50
}
