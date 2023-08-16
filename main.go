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
	err = cpu.loadMemory("Output/memoryFile.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = cpu.loadInputOuput("ioFile.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	cpu.print()
	cpu.run()
    cpu.printInputOuput()
}
