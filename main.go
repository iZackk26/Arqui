package main

import (
	"fmt"
)

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
	err = cpu.loadInputOuput("input.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	cpu.print()
	cpu.run()

}
