package main

import (
	"fmt"
	"genetic_pcb/genetic"
	"genetic_pcb/pcb"
	"image/color"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	fmt.Println("Hi!")

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	// fakeEdgeMovingComponent := &pcb.Component{
	// 	Nodes: []pcb.ComponentNode{{DX: 0, DY: 0}},
	// 	X1:    -5,
	// 	Y1:    -5,
	// 	X2:    5,
	// 	Y2:    5,
	// 	Kind:  pcb.EDGE_BREAKER_COMPONENT,
	// }

	netColors := []color.Color{
		&color.RGBA{238, 79, 121, 255},  // Rosa
		&color.RGBA{75, 139, 190, 255},  // Blu
		&color.RGBA{87, 178, 158, 255},  // Verde acqua
		&color.RGBA{249, 147, 79, 255},  // Arancione
		&color.RGBA{193, 109, 161, 255}, // Viola
		&color.RGBA{237, 201, 72, 255},  // Giallo
		&color.RGBA{145, 204, 211, 255}, // Celeste
		&color.RGBA{106, 141, 59, 255},  // Verde
		&color.RGBA{209, 86, 66, 255},   // Rosso
		&color.RGBA{155, 97, 64, 255},   // Marrone
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	// p1 := pcb.GeneratePcbFull(componentTemplates, 20, 6, maxX, maxY, randomGenerator)
	// p1 := pcb.GeneratePcbFull(componentTemplates, 7, 3, maxX, maxY, randomGenerator)
	p1 := pcb.GeneratePcbFull(componentTemplates, 25, 10, maxX, maxY, randomGenerator)
	p2 := pcb.ScrumblePcb(p1, maxX, maxY)
	pgo := pcb.NewPcbGeneticOperators(
		1,
		0.2,
		0.1,
		maxX,
		maxY,
		nodeSz,
		edgeSz,
		50,
		pcb.MutationParams{
			GlobalMutationWeight:                  10,
			RegenerateNetMutationWeight:           10,
			TranslateComponentGroupMutationWeight: 10,
			RotateComponentMutationWeight:         10,
			RerouteEdgeMutationWeight:             10,
			ChangePlaneMutationWeight:             10,
		},
		pcb.EvaluationParams{
			SamePlaneIntersectionCost:      1.0,
			DifferentPlaneIntersectionCost: 0.9,
			EdgeLengthCost:                 0.01,
			NonZeroPlaneEdgeCost:           0.09,
			OutOfBoundsCost:                100,
			MinDist:                        2,
		},
	)
	ctx := genetic.NewGeneticContext()
	c := pgo.CrossOver(p1, p2, ctx)
	fmt.Printf("%+v\n", p1.Genome)

	p1.ComputeGeometry(nodeSz, edgeSz)
	p2.ComputeGeometry(nodeSz, edgeSz)
	c.ComputeGeometry(nodeSz, edgeSz)

	pgo.Evaluate(c, ctx)

	c.Genome.Edges[0].Plane = 1

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
		0.01,
	)

	pcb.DrawPcbToImage(ga.CurrentPop[0].Individual, "first.png", int(maxX), int(maxY), 1, 1, netColors)

	prevValue := 0.0

	for i := 0; i < 20000; i++ {
		start := time.Now()

		ga.ComputeNextGeneration()

		elapsed := time.Since(start)

		fmt.Printf("Generation time: %s\n", elapsed)

		currValue := ga.CurrentPop[0].Fitness

		if currValue != prevValue {
			pcb.DrawPcbToImage(ga.CurrentPop[0].Individual, "best_nw.png", int(maxX), int(maxY), 1, 1, netColors)
			os.Rename("best_nw.png", "best.png")

			prevValue = currValue
		}

		fmt.Println(i)

		fmt.Printf("%v %v\n", ga.CurrentPop[0].Fitness, ga.CurrentPop[len(ga.CurrentPop)-1].Fitness)
		fmt.Println(ga.CurrentPop[0].Individual.Genome.Edges)

		// if ga.CurrentPop[0].Fitness > -1 {
		// 	break
		// }

	}
}
