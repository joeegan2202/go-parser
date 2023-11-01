package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	CONCAT = 0
	OR     = 1
	KLEENE = 2
)

var Epsilon = "\000"

type Transition struct {
	digest *string
	next   *Transition
}

type NFA struct {
	states       []string
	transitions  [][]*Transition
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

func (digraph Digraph) NodeList() []*DigraphNode {
	seen := make([]*DigraphNode, 0, 128)

	for _, node := range digraph.startingNodes {
		unique := true

		for _, seenNode := range seen {
			if node == seenNode {
				unique = false
				break
			}
		}

		if unique {
			seen = append(seen, node)
			node.appendChildren(&seen)
		}
	}

	return seen
}

func (node DigraphNode) appendChildren(seen *[]*DigraphNode) {
	for _, child := range node.next {
		unique := true

		for _, seenNode := range *seen {
			if child.node == seenNode {
				unique = false
				break
			}
		}

		if unique {
			*seen = append(*seen, child.node)
			child.node.appendChildren(seen)
		}
	}
}

func NormalizeRegex(regex string) *Digraph {
	digraph := new(Digraph)
	finalNode := new(DigraphNode)
	finalNode.name = regex
	initialNode := new(DigraphNode)
	initialNode.name = Epsilon
	initialNode.next = []struct {
		name string
		node *DigraphNode
	}{{regex, finalNode}}
	digraph.startingNodes = []*DigraphNode{initialNode}
	digraph.endingNodes = []*DigraphNode{finalNode}

	initialNode.normalize()

	return digraph
}

func (node DigraphNode) normalize() {
	for i, next := range node.next {
		if len(next.name) != 1 {

		}
	}
}

func nextOperation(regex string) (string, int, string) {
	first := regex[0]
	if first == '(' {

	}
}

func (digraph Digraph) ToNFA() *NFA {
	nfa := new(NFA)
	graphNodes := digraph.NodeList()
	nfa.states = make([]string, len(graphNodes)+1)
	nfa.states[0] = "initial"
	mapNodeToState := make(map[string]int, len(graphNodes))
	mapNodeToState["initial"] = 0

	for n, node := range graphNodes {
		mapNodeToState[node.name] = n + 1
		nfa.states[n+1] = node.name
	}

	nfa.transitions = make([][]*Transition, len(graphNodes)+1)
	for i := range nfa.transitions {
		nfa.transitions[i] = make([]*Transition, len(graphNodes)+1)
	}

	for _, node := range digraph.startingNodes {
		epsilonTransition := new(Transition)
		epsilonTransition.digest = &Epsilon
		epsilonTransition.next = nil
		nfa.transitions[0][mapNodeToState[node.name]] = epsilonTransition
	}

	for _, node := range graphNodes {
		for _, child := range node.next {
			newTransition := new(Transition)
			*newTransition = Transition{&child.name,
				nfa.transitions[mapNodeToState[node.name]][mapNodeToState[child.node.name]]}
			nfa.transitions[mapNodeToState[node.name]][mapNodeToState[child.node.name]] = newTransition
		}
	}

	nfa.finalStates = make([]int, len(digraph.endingNodes))
	for i, node := range digraph.endingNodes {
		nfa.finalStates[i] = mapNodeToState[node.name]
	}

	return nfa
}

func NewNFAFromRegex(regex string) *NFA {
	nfa := new(NFA)
	nfa.states = make([]string, 2, len(regex))
	nfa.states[0] = Epsilon
	nfa.states[1] = regex

	nfa.transitions = make([][]*Transition, len(nfa.states))
	for i := range nfa.states {
		nfa.transitions[i] = make([]*Transition, len(nfa.states))
	}
	newTransition := new(Transition)
	*newTransition = Transition{&regex, nil}
	nfa.transitions[0][1] = newTransition

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
			transition := nfa.transitions[i][t]
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
