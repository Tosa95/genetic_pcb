package main

import (
	"fmt"
	"genetic_pcb/genetic"
	"genetic_pcb/pcb"
	"image/color"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hi!")

	maxX, maxY := 500.0, 500.0
	nodeSz, edgeSz := 10.0, 5.0
	N := 1000

	componentTemplates := []pcb.Component{
		// Resistor
		{
			Nodes: []pcb.ComponentNode{{DX: -15, DY: 0}, {DX: 15, DY: 0}},
			X1:    -25,
			Y1:    -10,
			X2:    25,
			Y2:    10,
		},
		// Transistor
		{
			Nodes: []pcb.ComponentNode{{DX: -30, DY: 0}, {DX: 0, DY: 0}, {DX: 30, DY: 0}},
			X1:    -40,
			Y1:    -10,
			X2:    40,
			Y2:    10,
		},
	}

	netColors := []color.Color{
		&color.RGBA{0, 255, 0, 255},
		&color.RGBA{0, 0, 255, 255},
		&color.RGBA{0, 128, 128, 255},
		&color.RGBA{100, 20, 128, 255},
		&color.RGBA{200, 100, 100, 255},
		&color.RGBA{30, 100, 150, 255},
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	p1 := pcb.GeneratePcbFull(componentTemplates, 10, 3, maxX, maxY, randomGenerator)
	p2 := pcb.ScrumblePcb(p1, maxX, maxY)
	pgo := pcb.NewPcbGeneticOperators(
		1,
		0.2,
		0.1,
		1,
		maxX,
		maxY,
		nodeSz,
		edgeSz,
		50,
		pcb.MutationWeights{
			GlobalMutationWeight:                  1,
			RegenerateNetMutationWeight:           1,
			TranslateComponentGroupMutationWeight: 1,
			RotateComponentMutationWeight:         1,
		},
		0.01, // was 0.01
	)
	ctx := genetic.NewGeneticContext()
	c := pgo.CrossOver(p1, p2, ctx)
	fmt.Printf("%+v\n", p1.Genome)

	p1.ComputeGeometry(nodeSz, edgeSz)
	p2.ComputeGeometry(nodeSz, edgeSz)
	c.ComputeGeometry(nodeSz, edgeSz)

	pgo.Evaluate(c, ctx)

	pcb.DrawPcbToImage(p1, "p1.png", int(maxX), int(maxY), 1, 1, netColors)
	pcb.DrawPcbToImage(p2, "p2.png", int(maxX), int(maxY), 1, 1, netColors)
	pcb.DrawPcbToImage(c, "c.png", int(maxX), int(maxY), 1, 1, netColors)

	// // // p := pcb.GeneratePcb(25, 40, maxX, maxY)
	// // p := pcb.GeneratePcbWithNets(5, 6, maxX, maxY)
	// N := 1000
	// initialPop := make([]*pcb.Pcb, N)

	initialPop := make([]*pcb.Pcb, N)

	for i := 0; i < N; i++ {
		initialPop[i] = pcb.ScrumblePcb(p1, maxX, maxY)
	}

	ga := genetic.NewGeneticAlgorithm[*pcb.Pcb](
		initialPop,
		10,
		0.1,
		pgo,
		10,
		0.1,
	)

	pcb.DrawPcbToImage(ga.CurrentPop[0].Individual, "first.png", int(maxX), int(maxY), 1, 1, netColors)

	for i := 0; i < 5000; i++ {
		ga.ComputeNextGeneration()

		pcb.DrawPcbToImage(ga.CurrentPop[0].Individual, "best.png", int(maxX), int(maxY), 1, 1, netColors)

		fmt.Println(i)

		fmt.Printf("%v %v\n", ga.CurrentPop[0].Fitness, ga.CurrentPop[len(ga.CurrentPop)-1].Fitness)
		fmt.Println(ga.CurrentPop[0].Individual.Genome.Edges)

		// if ga.CurrentPop[0].Fitness > -1 {
		// 	break
		// }

	}
}
