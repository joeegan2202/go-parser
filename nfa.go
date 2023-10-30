package main

import "fmt"

type Transition struct {
	fromState int
	digest    *string
	toState   int
}

type NFA struct {
	states       []string
	transitions  []Transition
	initialState int
	finalStates  []int
}

func NewNFAFromRegex(regex string) *NFA {
	nfa := new(NFA)
	nfa.states = make([]string, 2, len(regex))
	nfa.states[0] = "\000"
	nfa.states[1] = regex

	nfa.transitions = make([]Transition, 1, len(regex))
	nfa.transitions[0] = Transition{0, &regex, 1}

	nfa.initialState = 0
	nfa.finalStates = make([]int, 1)
	nfa.finalStates[0] = 1

	return nfa
}

func (nfa *NFA) Print() {
	fmt.Printf("States:\nNumber\tName\n")
	for i, state := range nfa.states {
		fmt.Printf("%d\t%s\n", i, state)
	}

	fmt.Printf("\nTransitions:\n\t\t")
}
