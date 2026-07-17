package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	manager := NewCanvasManager()

	fmt.Println(
		"Comandos: c x y | l x y n | r a b c d | f x y r | q",
	)

	for scanner.Scan() {
		command, err := ParseCommand(scanner.Text())
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		res, err := command.Execute(manager)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if res.Quit {
			return
		}

		if res.Canvas != nil {
			fmt.Print(res.Canvas.String())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error leyendo la entrada:", err)
	}
}
