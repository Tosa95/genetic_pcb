package main

import (
	"fmt"
	"genetic_pcb/genetic"
	"genetic_pcb/pcb"
	"image/color"
)

func main() {
	fmt.Println("Hi!")

	maxX, maxY := 500.0, 500.0
	nodeSz, edgeSz := 10.0, 5.0

	// p := pcb.GeneratePcb(25, 40, maxX, maxY)
	p := pcb.GeneratePcbWithNets(5, 6, maxX, maxY)
	N := 1000
	initialPop := make([]*pcb.Pcb, N)

	p.Genome.Components = []pcb.Component{
		{
			Nodes: []pcb.ComponentNode{{0, 9, 0}},
			X1:    -30,
			Y1:    -30,
			X2:    30,
			Y2:    30,
			CX:    100,
			CY:    100,
		},
	}

	for i := 0; i < N; i++ {
		initialPop[i] = pcb.ScrumblePcb(p, maxX, maxY)
	}

	netColors := []color.Color{
		&color.RGBA{0, 255, 0, 255},
		&color.RGBA{0, 0, 255, 255},
		&color.RGBA{0, 128, 128, 255},
		&color.RGBA{100, 20, 128, 255},
		&color.RGBA{200, 100, 100, 255},
		&color.RGBA{30, 100, 150, 255},
	}

	ga := genetic.NewGeneticAlgorithm[*pcb.Pcb](
		initialPop,
		10,
		0.1,
		pcb.NewPcbGeneticOperators(
			1,
			0.05,
			0.5,
			10,
			maxX,
			maxY,
			nodeSz,
			edgeSz,
			50,
			0.5,
			0.1,
			0.1,
		),
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
