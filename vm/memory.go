package vm

type Memory struct {
	memory []int
}

func (mem *Memory) Store(offset int, val int) {

	mem.memory[offset] = val
}

func (mem *Memory) Load(offset int) int {
	return mem.memory[offset]
}
