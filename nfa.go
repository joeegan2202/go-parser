package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Transition struct {
	digest *string
	next   *Transition
}

type NFA struct {
	states       []string
	transitions  [][]Transition
	initialState int
	finalStates  []int
}

type DigraphNode struct {
	name string
	next []struct {
		name string
		node *DigraphNode
	}
}

type Digraph struct {
	startingNodes []*DigraphNode
	endingNodes   []*DigraphNode
}

func (digraph Digraph) CountNodes() int {
	accumulator := 0

	seen := make([]*DigraphNode, 128)

	for _, node := range digraph.startingNodes {
		for node != nil {
			unique := true

			for _, seenNode := range seen {
				if node == seenNode {
					unique = false
				}
			}

			if unique {
				accumulator++
				seen = append(seen, node)
			}
		}
	}

	return accumulator
}

func (digraph Digraph) ToNFA() *NFA {
	nfa := new(NFA)

	return nfa
}

func NewNFAFromRegex(regex string) *NFA {
	nfa := new(NFA)
	nfa.states = make([]string, 2, len(regex))
	nfa.states[0] = "\000"
	nfa.states[1] = regex

	nfa.transitions = make([][]Transition, len(nfa.states))
	for i := range nfa.states {
		nfa.transitions[i] = make([]Transition, len(nfa.states))
	}
	nfa.transitions[0][1] = Transition{&regex, nil}

	nfa.initialState = 0
	nfa.finalStates = make([]int, 1)
	nfa.finalStates[0] = 1

	return nfa
}

func (nfa *NFA) Print() {
	stateTable := table.NewWriter()
	stateTable.SetOutputMirror(os.Stdout)
	stateTable.AppendHeader(table.Row{"Number", "Name"})
	transitionHeader := make([]interface{}, len(nfa.states)+1)
	transitionHeader[0] = ""
	for i, state := range nfa.states {
		transitionHeader[i+1] = i
		stateTable.AppendRow([]interface{}{i, state})
	}
	fmt.Printf("States:\n")
	stateTable.Render()

	transitionTable := table.NewWriter()
	transitionTable.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	transitionTable.SetOutputMirror(os.Stdout)
	transitionTable.AppendHeader(transitionHeader)

	for i := 0; i < len(nfa.states); i++ {
		row := make(table.Row, len(nfa.states)+1)
		row[0] = i

		for t := range nfa.transitions[i] {
			transition := &nfa.transitions[i][t]
			s := ""
			if transition.digest != nil {
				s = *transition.digest
			}
			transition = transition.next

			for transition != nil {
				s += "\n" + *transition.digest
				transition = transition.next
			}

			row[t+1] = s
		}

		transitionTable.AppendRow(row)
	}

	fmt.Printf("\nTransitions:\n")
	transitionTable.Render()
}
