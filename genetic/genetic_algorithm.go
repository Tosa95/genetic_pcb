package genetic

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mroth/weightedrand/v2"
	"gonum.org/v1/gonum/floats"
)

type IndividualWithFitness[T any] struct {
	Individual T
	Fitness    float64
}

type GeneticAlgorithm[T fmt.Stringer] struct {
	CurrentPop                       []IndividualWithFitness[T]
	elitarismKeepN                   int
	pOfSelectingSecondParentRandomly float64
	geneticOperators                 GeneticOperators[T]
	parallelism                      int
	selfReproductionProb             float64
	randomGenerator                  *rand.Rand
}

func NewGeneticAlgorithm[T fmt.Stringer](initialPop []T, elitarismKeepN int, pOfSelectingSecondParentRandomly float64, geneticOperators GeneticOperators[T], parallelism int, selfReproductionProb float64) *GeneticAlgorithm[T] {
	evaluedPop := make([]IndividualWithFitness[T], len(initialPop))

	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	c := NewGeneticContext()

	for i, ind := range initialPop {
		geneticOperators.Grow(ind, c)
		eval := geneticOperators.Evaluate(ind, c)

		evaluedPop[i] = IndividualWithFitness[T]{Individual: ind, Fitness: eval}
	}

	ga := GeneticAlgorithm[T]{
		CurrentPop:                       evaluedPop,
		elitarismKeepN:                   elitarismKeepN,
		pOfSelectingSecondParentRandomly: pOfSelectingSecondParentRandomly,
		geneticOperators:                 geneticOperators,
		parallelism:                      parallelism,
		selfReproductionProb:             selfReproductionProb,
		randomGenerator:                  randomGenerator,
	}

	ga.sortPop()

	return &ga
}

func (ga *GeneticAlgorithm[T]) sortPop() {
	sort.SliceStable(ga.CurrentPop, func(i, j int) bool { return ga.CurrentPop[i].Fitness > ga.CurrentPop[j].Fitness })
}

// TODO: use floats instead of ints...
func popToChoices(eval []float64) []weightedrand.Choice[int, int] {
	min := floats.Min(eval)
	max := floats.Max(eval)

	res := make([]weightedrand.Choice[int, int], len(eval))

	for i, e := range eval {
		if min < max {
			res[i] = weightedrand.NewChoice(i, int((e-min)*10000))
		} else {
			res[i] = weightedrand.NewChoice(i, 1)
		}

	}

	return res
}

func (ga *GeneticAlgorithm[T]) generateChild(p1 int, p2 int, c *GeneticContext) IndividualWithFitness[T] {
	// fmt.Println(ga.CurrentPop[p1].Individual.String())
	child := ga.geneticOperators.CrossOver(ga.CurrentPop[p1].Individual, ga.CurrentPop[p2].Individual, c)
	// fmt.Println(ga.CurrentPop[p1].Individual.String())
	ga.geneticOperators.Mutate(child, c)
	// fmt.Println(ga.CurrentPop[p1].Individual.String())
	ga.geneticOperators.Grow(child, c)
	// fmt.Println(ga.CurrentPop[p1].Individual.String())
	fitness := ga.geneticOperators.Evaluate(child, c)

	// fmt.Println(ga.CurrentPop[p1].Individual.String())
	// fmt.Println()

	return IndividualWithFitness[T]{
		Individual: child,
		Fitness:    fitness,
	}
}

func (ga *GeneticAlgorithm[T]) generateChilds(inputChannel chan []int, outputChannel chan IndividualWithFitness[T], wg *sync.WaitGroup) {
	defer wg.Done()

	c := NewGeneticContext()

	for parents := range inputChannel {
		p1, p2 := parents[0], parents[1]
		outputChannel <- ga.generateChild(p1, p2, c)

		// if c.RandomGenerator.Float64() < ga.selfReproductionProb {
		// 	ga.geneticOperators.Mutate(ga.CurrentPop[p1].Individual, c)
		// }

	}
}

func (ga *GeneticAlgorithm[T]) ComputeNextGeneration() {

	toBeKept := ga.CurrentPop[:ga.elitarismKeepN]

	toBeGenerated := len(ga.CurrentPop) - len(toBeKept)

	generated := make([]IndividualWithFitness[T], toBeGenerated)
	evals := make([]float64, len(ga.CurrentPop))

	for i, f := range ga.CurrentPop {
		evals[i] = f.Fitness
	}

	inputChan := make(chan []int, toBeGenerated)
	outputChan := make(chan IndividualWithFitness[T], toBeGenerated)
	chooser, _ := weightedrand.NewChooser(popToChoices(evals)...)

	for i := 0; i < toBeGenerated; i++ {

		if ga.randomGenerator.Float64() < ga.pOfSelectingSecondParentRandomly {
			inputChan <- []int{chooser.Pick(), ga.randomGenerator.Intn(len(ga.CurrentPop))}
		} else {
			inputChan <- []int{chooser.Pick(), chooser.Pick()}
		}

	}

	close(inputChan)

	wg := sync.WaitGroup{}

	for i := 0; i < ga.parallelism; i++ {
		wg.Add(1)
		go ga.generateChilds(inputChan, outputChan, &wg)
	}

	wg.Wait()

	close(outputChan)

	for i := 0; i < toBeGenerated; i++ {
		c := <-outputChan

		generated[i] = c
	}

	ga.CurrentPop = append(toBeKept, generated...)

	ga.sortPop()

}
