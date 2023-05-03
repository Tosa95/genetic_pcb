package pcb_test

import (
	"genetic_pcb/pcb"
	"testing"
)

func runTest(t *testing.T, g pcb.Genome, expected float64) {
	p := pcb.NewPcb(&g)

	p.ComputeGeometry(10, 5)

	// res := pcb.EvaluatePcb(p, 10, 0)

	// if res != expected {
	// 	t.Errorf("Expected %v, got %v", expected, res)
	// }
}

func TestEvaluatePcb1(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10, 0},
				{50, 11, 0},
				{100, 100, 0},
				{83, 12, 0},
			},
			Edges: []pcb.Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 2, To: 3},
			},
		},
		0.0,
	)
}

func TestEvaluatePcb2(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10, 0},
				{50, 11, 0},
				{100, 100, 0},
				{63, 12, 0},
			},
			Edges: []pcb.Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 2, To: 3},
			},
		},
		-4.0,
	)
}

func TestEvaluatePcb3(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10, 0},
				{50, 11, 0},
				{100, 100, 0},
				{103, 12, 0},
				{22, 72, 0},
			},
			Edges: []pcb.Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 3, To: 4},
			},
		},
		-1.0,
	)
}

func TestEvaluatePcb4(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10, 0},
				{50, 11, 0},
				{100, 100, 0},
				{194, 12, 0},
				{82, 172, 0},
			},
			Edges: []pcb.Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 3, To: 4},
			},
		},
		0.0,
	)
}

func TestEvaluatePcb5(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10, 0},
				{50, 11, 0},
				{100, 100, 0},
				{194, 12, 0},
				{52, 172, 0},
			},
			Edges: []pcb.Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 3, To: 4},
			},
		},
		-2.0,
	)
}
