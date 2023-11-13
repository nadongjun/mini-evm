package vm

import "fmt"

type Memory struct {
	memory []int
}

func (mem *Memory) Store(offset int, val int) {
	fmt.Println("val", offset, val)
	mem.memory = append(mem.memory, val)
}

func (mem *Memory) Load(offset int) int {
	return mem.memory[offset]
}

func (mem *Memory) LoadRange(length int, offset int) []int {

	if offset < 0 {
		panic(fmt.Sprintf("InvalidMemoryAccess %v is already registered."))

	}
	data := make([]int, length)
	for i := 0; i < length; i++ {
		data[i] = mem.Load(offset + i)
	}
	return data
}
