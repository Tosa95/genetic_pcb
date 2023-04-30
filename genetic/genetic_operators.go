package genetic

import (
	"math/rand"
	"time"
)

type GeneticContext struct {
	RandomGenerator *rand.Rand
}

func NewGeneticContext() *GeneticContext {
	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	c := GeneticContext{
		RandomGenerator: randomGenerator,
	}

	return &c
}

type IndividualEvaluator[T any] interface {
	Evaluate(i T, c *GeneticContext) float64
}

type CrossoverManager[T any] interface {
	CrossOver(i1 T, i2 T, c *GeneticContext) T
}

type MutationManager[T any] interface {
	Mutate(i T, c *GeneticContext)
}

type GrowthManager[T any] interface {
	Grow(i T, c *GeneticContext)
}

type GeneticOperators[T any] interface {
	IndividualEvaluator[T]
	CrossoverManager[T]
	MutationManager[T]
	GrowthManager[T]
}
