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
		&color.RGBA{0x1F, 0x77, 0xB4, 0xFF},
		&color.RGBA{0xFF, 0x7F, 0x0E, 0xFF},
		&color.RGBA{0x2C, 0xA0, 0x2C, 0xFF},
		&color.RGBA{0xD6, 0x27, 0x28, 0xFF},
		&color.RGBA{0x94, 0x67, 0xBD, 0xFF},
		&color.RGBA{0x8C, 0x56, 0x46, 0xFF},
		&color.RGBA{0xE3, 0x77, 0xC2, 0xFF},
		&color.RGBA{0x7F, 0x7F, 0x7F, 0xFF},
		&color.RGBA{0xBC, 0xBD, 0x22, 0xFF},
		&color.RGBA{0x17, 0xBE, 0xCF, 0xFF},
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	p1 := pcb.GeneratePcbFull(componentTemplates, 20, 6, maxX, maxY, randomGenerator)
	// p1 := pcb.GeneratePcbFull(componentTemplates, 7, 3, maxX, maxY, randomGenerator)
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
			RerouteEdgeMutationWeight:             1,
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

	// pgo.MoveEdgeMutation(c, ctx)
	// c.ComputeGeometry(nodeSz, edgeSz)

	// pcb.DrawPcbToImage(c, "c2.png", int(maxX), int(maxY), 1, 1, netColors)

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
