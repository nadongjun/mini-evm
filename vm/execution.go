package vm

import (
	"encoding/binary"
	"fmt"
)

// Execution layer
type ExecutionContext struct {
	code       []byte
	stack      Stack
	memory     Memory
	pc         int
	stopped    bool
	returndata []int
}

// Evm instruction
type Instruction struct {
	opcode  int
	name    string
	execute func(*ExecutionContext)
}

var Instructions []Instruction
var InstructionsByOpcode = make(map[int]Instruction)
var (
	STOP  = RegisterInstruction(0x00, "STOP", func(ctx *ExecutionContext) { ctx.Stop() })
	PUSH1 = RegisterInstruction(0x60, "PUSH1", func(ctx *ExecutionContext) { ctx.stack.push(ctx.ReadCode(1)) })
	ADD   = RegisterInstruction(0x01, "ADD", func(ctx *ExecutionContext) { ctx.stack.push((ctx.stack.pop() + ctx.stack.pop()) % 256) })
	MUL   = RegisterInstruction(0x02, "MUL", func(ctx *ExecutionContext) { ctx.stack.push((ctx.stack.pop() * ctx.stack.pop()) % 256) })
	MLOAD = RegisterInstruction(
		0x51,
		"MLOAD",
		func(ctx *ExecutionContext) {
			temp := intsToBytes(ctx.memory.LoadRange(ctx.stack.pop(), 32))
			// Word to int
			ctx.stack.push(int(binary.BigEndian.Uint64(temp)))
		},
	)
	MSTORE8 = RegisterInstruction(
		0x53,
		"MSTORE8",
		func(ctx *ExecutionContext) {
			address := ctx.stack.pop()
			value := ctx.stack.pop() % 256

			ctx.memory.Store(address, value)
		},
	)
	RETURN = RegisterInstruction(
		0xf3,
		"RETURN",
		func(ctx *ExecutionContext) {
			dataSize := ctx.stack.pop()
			dataOffset := ctx.stack.pop()
			ctx.SetReturnData(dataOffset, dataSize)
		},
	)
	JUMP = RegisterInstruction(
		0x56,
		"JUMP",
		func(ctx *ExecutionContext) {
			ctx.SetProgramCounter(ctx.stack.pop())
		},
	)
	JUMPI = RegisterInstruction(
		0x57,
		"JUMPI",
		func(ctx *ExecutionContext) {
			target_pc, cond := ctx.stack.pop(), ctx.stack.pop()
			if cond != 0 {
				ctx.SetProgramCounter(target_pc)
			}
		},
	)
	PC = RegisterInstruction(
		0x58,
		"PC",
		func(ctx *ExecutionContext) {
			ctx.stack.push(ctx.pc)
		},
	)
)

func NewExecutionContext(code []byte, pc int, stack Stack, memory Memory) *ExecutionContext {
	return &ExecutionContext{
		code:    code,
		pc:      pc,
		stack:   stack,
		memory:  memory,
		stopped: false,
	}
}

func (exe *ExecutionContext) Stop() {
	exe.stopped = true
}

func (exe *ExecutionContext) ReadCode(numBytes int) int {
	fmt.Println(int(exe.code[exe.pc : exe.pc+numBytes][0]))
	value := int(exe.code[exe.pc : exe.pc+numBytes][0])
	exe.pc += numBytes

	return value
}

func (exe *ExecutionContext) SetReturnData(offset, length int) {
	exe.stopped = true
	exe.returndata = exe.memory.LoadRange(offset, length)
}

func (exe *ExecutionContext) SetProgramCounter(pc int) {
	exe.pc = pc
}

func RegisterInstruction(opcode int, name string, executeFunc func(*ExecutionContext)) *Instruction {
	instruction := &Instruction{opcode: opcode, name: name, execute: executeFunc}
	Instructions = append(Instructions, *instruction)

	if _, exists := InstructionsByOpcode[opcode]; exists {
		panic(fmt.Sprintf("Opcode %v is already registered.", opcode))
	}
	InstructionsByOpcode[opcode] = *instruction

	return instruction
}

func DecodeOpcode(context *ExecutionContext) (Instruction, error) {
	if context.pc < 0 || context.pc >= len(context.code) {
		return Instruction{}, fmt.Errorf("InvalidCodeOffset: code=%v, pc=%v", context.code, context.pc)
	}

	opcode := context.ReadCode(1)
	instruction, exists := InstructionsByOpcode[opcode]
	fmt.Println(instruction.name)
	if !exists {
		return Instruction{}, fmt.Errorf("UnknownOpcode: opcode=%v", opcode)
	}

	return instruction, nil
}

func Run(code []byte) int {
	// Executes code in a fresh context.
	context := NewExecutionContext(code, 0, Stack{}, Memory{})

	for !context.stopped {
		instruction, err := DecodeOpcode(context)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		instruction.execute(context)

		fmt.Println("pc", context.pc)
		fmt.Println("stack", context.stack.stack)
		fmt.Println("memeory", context.memory.memory)
		fmt.Println()
	}
	fmt.Println("return data= ", context.returndata[0])
	return context.returndata[0]
}

// Util
// [1, 2, 3, 4] -> []byte
func intsToBytes(ints []int) []byte {
	byteSlice := make([]byte, len(ints)*4)

	for i, v := range ints {
		binary.BigEndian.PutUint32(byteSlice[i*4:], uint32(v))
	}

	return byteSlice
}
