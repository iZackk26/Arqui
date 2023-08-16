package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type CPU struct {
	accumulator  int64
	instructions *bytes.Buffer
	memory       [2048]int64
	inputOuput   [10]int64
}

func (c *CPU) loadInstructions(filename string) error {
	instructions := &bytes.Buffer{}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Fields(line)
		for _, token := range temp {
			num, err := strconv.ParseInt(token, 2, 64)
			if err != nil {
				return err
			}

			binary.Write(instructions, binary.LittleEndian, num)
		}
	}
	c.instructions = instructions
	return nil
}

func (c *CPU) dumpMemory(filename string) {
	file, err := os.Create(filename) // If file already exists, it will be truncated
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// for i := 0; i < len(c.memory); i++ {
	// 	if c.memory[i] != 0 {
	// 		fmt.Fprintf(file, "%b %b\n", i, c.memory[i]) // B means binary
	// 	}
	// }

	for idx, mem := range c.memory {
		if mem != 0 {
			fmt.Fprintf(file, "%b %b\n", idx, mem)
		}
	}
}
func (c *CPU) dumpMemoryIO(filename string) {
	// file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for _, mem := range c.inputOuput {
		fmt.Fprintf(file, "%b\n", mem)
	}
	fmt.Fprintf(file, "\n")
}

func (c *CPU) loadMemory(filename string) error {
	var memory [2048]int64
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Fields(line)

		address, err := strconv.ParseInt(temp[0], 2, 64)
		if err != nil {
			return err
		}
		value, err := strconv.ParseInt(temp[1], 2, 64)
		if err != nil {
			return err
		}
		memory[address] = value
	}
	c.memory = memory
	return nil
}

func (c *CPU) loadInputOuput(filename string) error {
	var inputOuput [10]int64
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	text, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Fields(string(text))

	for idx, line := range lines {
		num, err := strconv.ParseInt(line, 2, 64)
		if err != nil {
			return err
		}
		inputOuput[idx] = num
	}

	c.inputOuput = inputOuput
	return nil
}

func (c *CPU) fetch() int64 {
	var instruction int64
	binary.Read(c.instructions, binary.LittleEndian, &instruction)
	return instruction
}

func (c *CPU) printMemory() {
	for idx, mem := range c.memory {
		if mem != 0 {
			fmt.Printf("Memory %v: %v\n", idx,mem)
		}
	}
}

func (c *CPU) printInputOuput() {
    for idx, mem := range c.inputOuput {
        if mem != 0 {
            fmt.Printf("InputOuput %v: %v\n", idx,mem)
        }
    }
}

func (c *CPU) print() {
	fmt.Printf("Accumulator: %d\n", c.accumulator)
	// fmt.Printf("Instructions: %d\n", c.instructions)
	c.printMemory()
	fmt.Println()
}

func (c *CPU) run() {
	for {
		instruction := c.fetch()
		if instruction == 0 {
			break
		}
		err := c.execute(instruction)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(instruction)
		c.print()
		c.dumpMemory("Output/memoryFileOutput.txt")
	}
}

func (c *CPU) execute(instruction int64) error {
	switch instruction {
	case 1:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator = c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory) // skip empty bits
	case 2:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.memory[memory] = c.accumulator
		binary.Read(c.instructions, binary.LittleEndian, &memory) // skip empty bits
	case 3:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator += c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory)
	case 4:
		var memory1 int64
		var memory2 int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory1)
		if err != nil {
			return err
		}
		err = binary.Read(c.instructions, binary.LittleEndian, &memory2)
		if err != nil {
			return err
		}
		c.accumulator += c.memory[memory1] + c.memory[memory2]

	case 5:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator -= c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory)

	case 6:
		var memory1 int64
		var memory2 int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory1)
		if err != nil {
			return err
		}
		err = binary.Read(c.instructions, binary.LittleEndian, &memory2)
		if err != nil {
			return err
		}
		c.memory[memory2] = c.accumulator - c.memory[memory1]
	case 7:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator *= c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory)
	case 8:
		// Sum two memory locations and store in memory 1
		var memory1 int64
		var memory2 int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory1)
		if err != nil {
			return err
		}
		err = binary.Read(c.instructions, binary.LittleEndian, &memory2)
		if err != nil {
			return err
		}
		c.memory[memory1] += c.memory[memory2]
	case 9:
		// Multiply memory 1 by accumulator and store in memory 2
		var memory1 int64
		var memory2 int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory1)
		if err != nil {
			return err
		}
		err = binary.Read(c.instructions, binary.LittleEndian, &memory2)
		if err != nil {
			return err
		}
		c.memory[memory2] = c.memory[memory1] * c.accumulator
	case 10:
		// Divide accumulator by memory 1 and store it en accumulator
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator /= c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory)
	case 11:
		// Divide the accumulator by memory 1 and store it in memory 2
		var memory1 int64
		var memory2 int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory1)
		if err != nil {
			return err
		}
		err = binary.Read(c.instructions, binary.LittleEndian, &memory2)
		if err != nil {
			return err
		}
		c.memory[memory2] = c.accumulator / c.memory[memory1]
	case 12:
		// Loads The inputOuput and stores it in the accumulator
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.accumulator = c.inputOuput[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory)
		c.dumpMemoryIO("Output/outputFile.txt")
	case 13:
		// Loads the accumulator and stores it in the inputOuput
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
		c.inputOuput[memory] = c.accumulator
		binary.Read(c.instructions, binary.LittleEndian, &memory)
		c.dumpMemoryIO("Output/outputFile.txt")
	default:
		return fmt.Errorf("Unknown instruction: %d", instruction)

	}
	return nil

}
