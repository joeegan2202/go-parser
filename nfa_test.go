package main

import (
	"strings"
	"testing"
)

func TestNodelist(t *testing.T) {
	{ // Test a -> b -> c straight graph
		nodeC := DigraphNode{"c", make([]struct {
			name string
			node *DigraphNode
		}, 0)}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeC}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}}}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA}, []*DigraphNode{&nodeC}}
		got := testGraph1.NodeList()
		expected := []*DigraphNode{&nodeA, &nodeB, &nodeC}

		if len(got) != len(expected) {
			t.Errorf("Did not get correct number of nodes! Expected %d, got %d!", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("Unexpected node found! Expected {Node %s}, got {Node %s}!", expected[i].name, got[i].name)
			}
		}
	}

	{ // Test (a) -> (b, c) -> (d), diamond graph
		nodeD := DigraphNode{"d", make([]struct {
			name string
			node *DigraphNode
		}, 0)}
		nodeC := DigraphNode{"c", []struct {
			name string
			node *DigraphNode
		}{{"c", &nodeD}}}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeD}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}, {"a", &nodeC}}}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA}, []*DigraphNode{&nodeD}}
		got := testGraph1.NodeList()
		expected := []*DigraphNode{&nodeA, &nodeB, &nodeD, &nodeC}

		if len(got) != len(expected) {
			t.Errorf("Did not get correct number of nodes! Expected %d, got %d!", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("Unexpected node found! Expected {Node %s}, got {Node %s}!", expected[i].name, got[i].name)
			}
		}
	}

	{ // Test (a) -> (b) -> (c) -> (d) -> (b) -> (c) -> etc..., cyclic graph
		nodeD := DigraphNode{"d", make([]struct {
			name string
			node *DigraphNode
		}, 1)}
		nodeC := DigraphNode{"c", []struct {
			name string
			node *DigraphNode
		}{{"c", &nodeD}}}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeD}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}, {"a", &nodeC}}}
		nodeD.next[0] = struct {
			name string
			node *DigraphNode
		}{"d", &nodeB}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA}, []*DigraphNode{&nodeC}}
		got := testGraph1.NodeList()
		expected := []*DigraphNode{&nodeA, &nodeB, &nodeD, &nodeC}

		if len(got) != len(expected) {
			t.Errorf("Did not get correct number of nodes! Expected %d, got %d!", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("Unexpected node found! Expected {Node %s}, got {Node %s}!", expected[i].name, got[i].name)
			}
		}
	}

	{ // Test (a) -> (b) -> (c) -> (d) -> (b) -> (c) -> etc..., cyclic graph with a and b starting nodes
		nodeD := DigraphNode{"d", make([]struct {
			name string
			node *DigraphNode
		}, 1)}
		nodeC := DigraphNode{"c", []struct {
			name string
			node *DigraphNode
		}{{"c", &nodeD}}}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeD}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}, {"a", &nodeC}}}
		nodeD.next[0] = struct {
			name string
			node *DigraphNode
		}{"d", &nodeB}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA, &nodeB}, []*DigraphNode{&nodeC}}
		got := testGraph1.NodeList()
		expected := []*DigraphNode{&nodeA, &nodeB, &nodeD, &nodeC}

		if len(got) != len(expected) {
			t.Errorf("Did not get correct number of nodes! Expected %d, got %d!", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("Unexpected node found! Expected {Node %s}, got {Node %s}!", expected[i].name, got[i].name)
			}
		}
	}
}

func TestToNFA(t *testing.T) {
	{ // Test a -> b -> c straight graph
		nodeC := DigraphNode{"c", make([]struct {
			name string
			node *DigraphNode
		}, 0)}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeC}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}}}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA}, []*DigraphNode{&nodeC}}
		got := testGraph1.ToNFA()

		if strings.Compare(got.states[0], "initial") != 0 || strings.Compare(got.states[1], "a") != 0 || strings.Compare(got.states[2], "b") != 0 || strings.Compare(got.states[3], "c") != 0 {
			t.Errorf("Did not get states correctly!")
		}

		if got.transitions[0][0] != nil || got.transitions[0][1].digest != &Epsilon || got.transitions[0][2] != nil || got.transitions[0][3] != nil {
			t.Errorf("Did not get row 0 correctly!")
		}

		if got.transitions[1][0] != nil || got.transitions[1][1] != nil || strings.Compare(*got.transitions[1][2].digest, "a") != 0 || got.transitions[1][3] != nil {
			t.Errorf("Did not get row 1 correctly!")
		}

		if got.transitions[2][0] != nil || got.transitions[2][1] != nil || got.transitions[2][2] != nil || strings.Compare(*got.transitions[2][3].digest, "b") != 0 {
			t.Errorf("Did not get row 2 correctly!")
		}

		if got.finalStates[0] != 3 {
			t.Errorf("Did not get final states correctly!")
		}
	}

	{ // Test (a) -> (b, c) -> (d), diamond graph
		nodeD := DigraphNode{"d", make([]struct {
			name string
			node *DigraphNode
		}, 0)}
		nodeC := DigraphNode{"c", []struct {
			name string
			node *DigraphNode
		}{{"c", &nodeD}}}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeD}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}, {"a", &nodeC}}}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA}, []*DigraphNode{&nodeD}}
		got := testGraph1.ToNFA()

		if strings.Compare(got.states[0], "initial") != 0 || strings.Compare(got.states[1], "a") != 0 || strings.Compare(got.states[2], "b") != 0 || strings.Compare(got.states[3], "d") != 0 || strings.Compare(got.states[4], "c") != 0 {
			t.Errorf("Did not get states correctly!")
		}

		if got.transitions[0][0] != nil || got.transitions[0][1].digest != &Epsilon || got.transitions[0][2] != nil || got.transitions[0][3] != nil || got.transitions[0][4] != nil {
			t.Errorf("Did not get row 0 correctly!")
		}

		if got.transitions[1][0] != nil || got.transitions[1][1] != nil || strings.Compare(*got.transitions[1][2].digest, "a") != 0 || strings.Compare(*got.transitions[1][4].digest, "a") != 0 || got.transitions[1][3] != nil {
			t.Errorf("Did not get row 1 correctly!")
		}

		if got.transitions[2][0] != nil || got.transitions[2][1] != nil || got.transitions[2][2] != nil || strings.Compare(*got.transitions[2][3].digest, "b") != 0 || got.transitions[2][4] != nil {
			t.Errorf("Did not get row 2 correctly!")
		}

		if got.transitions[3][0] != nil || got.transitions[3][1] != nil || got.transitions[3][2] != nil || got.transitions[3][3] != nil || got.transitions[3][4] != nil {
			t.Errorf("Did not get row 3 correctly!")
		}

		if got.transitions[4][0] != nil || got.transitions[4][1] != nil || got.transitions[4][2] != nil || strings.Compare(*got.transitions[4][3].digest, "c") != 0 || got.transitions[4][4] != nil {
			t.Errorf("Did not get row 4 correctly!")
		}

		if got.finalStates[0] != 3 {
			t.Errorf("Did not get final states correctly!")
		}
	}

	{ // Test (a) -> (b) -> (c) -> (d) -> (b) -> (c) -> etc..., cyclic graph with a and b starting nodes
		nodeD := DigraphNode{"d", make([]struct {
			name string
			node *DigraphNode
		}, 1)}
		nodeC := DigraphNode{"c", []struct {
			name string
			node *DigraphNode
		}{{"c", &nodeD}}}
		nodeB := DigraphNode{"b", []struct {
			name string
			node *DigraphNode
		}{{"b", &nodeC}}}
		nodeA := DigraphNode{"a", []struct {
			name string
			node *DigraphNode
		}{{"a", &nodeB}}}
		nodeD.next[0] = struct {
			name string
			node *DigraphNode
		}{"d", &nodeB}
		testGraph1 := Digraph{[]*DigraphNode{&nodeA, &nodeB}, []*DigraphNode{&nodeC}}
		got := testGraph1.ToNFA()

		if strings.Compare(got.states[0], "initial") != 0 || strings.Compare(got.states[1], "a") != 0 || strings.Compare(got.states[2], "b") != 0 || strings.Compare(got.states[3], "c") != 0 || strings.Compare(got.states[4], "d") != 0 {
			t.Errorf("Did not get states correctly!")
		}

		if got.transitions[0][0] != nil || got.transitions[0][1].digest != &Epsilon || got.transitions[0][2].digest != &Epsilon || got.transitions[0][3] != nil || got.transitions[0][4] != nil {
			t.Errorf("Did not get row 0 correctly!")
		}

		if got.transitions[1][0] != nil || got.transitions[1][1] != nil || strings.Compare(*got.transitions[1][2].digest, "a") != 0 || got.transitions[1][3] != nil || got.transitions[1][4] != nil {
			t.Errorf("Did not get row 1 correctly!")
		}

		if got.transitions[2][0] != nil || got.transitions[2][1] != nil || got.transitions[2][2] != nil || strings.Compare(*got.transitions[2][3].digest, "b") != 0 || got.transitions[2][4] != nil {
			t.Errorf("Did not get row 2 correctly!")
		}

		if got.transitions[3][0] != nil || got.transitions[3][1] != nil || got.transitions[3][2] != nil || got.transitions[3][3] != nil || strings.Compare(*got.transitions[3][4].digest, "c") != 0 {
			t.Errorf("Did not get row 3 correctly!")
		}

		if got.transitions[4][0] != nil || got.transitions[4][1] != nil || strings.Compare(*got.transitions[4][2].digest, "d") != 0 || got.transitions[4][3] != nil || got.transitions[4][4] != nil {
			t.Errorf("Did not get row 4 correctly!")
		}

		if got.finalStates[0] != 3 {
			t.Errorf("Did not get final states correctly!")
		}
	}
}
