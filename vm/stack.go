package vm

import "fmt"

type Stack struct {
	stack    []int
	maxDepth int
}

func (s *Stack) push(val int) {
	// Stackoverflow
	if (len(s.stack) + 1) > s.maxDepth {
		fmt.Errorf("StackOverflow")
	}
	s.stack = append(s.stack, val)
}

func (s *Stack) pop() int {
	// Stackunderflow
	if len(s.stack) == 0 {
		fmt.Errorf("StackUnderflow")
	}
	length := len(*&s.stack)
	lastEle := (*&s.stack)[length-1]
	*&s.stack = (*&s.stack)[:length-1]
	return lastEle
}
