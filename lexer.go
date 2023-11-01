package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file! - ", err)
	}

	bufferedFile := bufio.NewReader(file)

	for byte, err := bufferedFile.ReadByte(); err == nil; byte, err = bufferedFile.ReadByte() {
		fmt.Printf("%c (0x%X)\n", byte, byte)
		empty := ""
		fmt.Scanln(empty)
	}

	fmt.Println("Finished scanning!")

	nfa := NewNFAFromRegex("ab*cdef(acd)*")
	bc := "bc"
	abcdefg := "abcdefg"
	nexttrans := Transition{&abcdefg, nil}
	newTransition := Transition{&bc, &nexttrans}
	nfa.transitions[1][0] = &newTransition

	nfa.Print()

	fmt.Println("Done printing! Goodbye!")
}
