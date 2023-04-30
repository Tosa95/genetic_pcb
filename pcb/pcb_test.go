package pcb_test

import (
	"genetic_pcb/pcb"
	"testing"
)

func runTest(t *testing.T, g pcb.Genome, expected float64) {
	p := pcb.NewPcb(&g)

	p.ComputeGeometry(10, 5)

	res := pcb.EvaluatePcb(p, 10)

	if res != expected {
		t.Errorf("Expected %v, got %v", expected, res)
	}
}

func TestEvaluatePcb1(t *testing.T) {
	runTest(
		t,
		pcb.Genome{
			Nodes: []pcb.Node{
				{10, 10},
				{50, 11},
				{100, 100},
				{83, 12},
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
				{10, 10},
				{50, 11},
				{100, 100},
				{63, 12},
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
				{10, 10},
				{50, 11},
				{100, 100},
				{103, 12},
				{22, 72},
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
				{10, 10},
				{50, 11},
				{100, 100},
				{194, 12},
				{82, 172},
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
				{10, 10},
				{50, 11},
				{100, 100},
				{194, 12},
				{52, 172},
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
