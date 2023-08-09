package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)
type CPU struct {
	accumulator  int64
	instructions *bytes.Buffer
	memory       [3000]int64
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

func (c *CPU)dumpMemory(filename string) {
    file, err := os.Create(filename) // If file already exists, it will be truncated
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()
    for i := 0; i < len(c.memory); i++ {
        if c.memory[i] != 0 {
            fmt.Fprintf(file, "%b %b\n", i, c.memory[i]) // B means binary
        }
    }
}

func (c *CPU) loadMemory(filename string) error {
	var memory [3000]int64
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

func (c *CPU) fetch() int64 {
	var instruction int64
	binary.Read(c.instructions, binary.LittleEndian, &instruction)
	return instruction
}

func (c *CPU) printMemory() {
	for i := 0; i < len(c.memory); i++ {
		if c.memory[i] != 0 {
			fmt.Printf("Memory[%d]: %d\n", i, c.memory[i])
		}
        // fmt.Printf("Memory[%d]: %d\n", i, c.memory[i]) 
	}
    // fmt.Printf("Memory[]: %v\n", c.memory) Mostrar esto
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
		fmt.Printf("Instruction: %d\n", instruction)
		c.print()
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
        fmt.Println("Memory: ", c.memory[memory])
		c.accumulator = c.memory[memory]
		binary.Read(c.instructions, binary.LittleEndian, &memory) // skip empty bits

	case 2:
        var memory int64 
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
        fmt.Println("Memory: ", memory, c.memory[memory])
		c.memory[memory] = c.accumulator
		binary.Read(c.instructions, binary.LittleEndian, &memory) // skip empty bits

	case 3:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
        fmt.Println("Memory: ", c.memory[memory])
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
        fmt.Println("Memory: ", c.memory[memory1], c.memory[memory2])
		c.accumulator += c.memory[memory1] + c.memory[memory2]

	case 5:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
        fmt.Println("Memory: ", c.memory[memory])
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
        fmt.Println("Memory: ", c.memory[memory1], c.memory[memory2])
		c.memory[memory2] = c.accumulator - c.memory[memory1]

	case 7:
		var memory int64
		err := binary.Read(c.instructions, binary.LittleEndian, &memory)
		if err != nil {
			return err
		}
        fmt.Println("Memory: ", c.memory[memory])
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
        fmt.Println("Memory: ", c.memory[memory1], c.memory[memory2])
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
        fmt.Println("Memory: ", c.memory[memory1], c.memory[memory2])
        c.memory[memory2] = c.memory[memory1] * c.accumulator
    case 10:
        // Divide accumulator by memory 1 and store it en accumulator
        var memory int64
        err := binary.Read(c.instructions, binary.LittleEndian, &memory)
        if err != nil {
            return err
        }
        fmt.Println("Memory: ", c.memory[memory])
        if c.memory[memory] == 0 {
            c.accumulator = 0
        } else {
        c.accumulator /= c.memory[memory]
        }
        
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
        fmt.Println("Memory: ", c.memory[memory1], c.memory[memory2])
        if c.memory[memory1] == 0 {
            c.memory[memory2] = 0
        } else {
        c.memory[memory2] = c.accumulator / c.memory[memory1]
    }
	default:
		return fmt.Errorf("Unknown instruction: %d", instruction)

	}
	return nil

}

func main() {
	var cpu CPU
	err := cpu.loadInstructions("instructions.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = cpu.loadMemory("memoryA.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	cpu.print()
	cpu.run()
    cpu.dumpMemory("Test.txt")

}
