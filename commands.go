package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CommandResult struct {
	Canvas *Canvas
	Quit   bool
}

type Command interface {
	Execute(manager *CanvasManager) (CommandResult, error)
}

type CreateCanvasCommand struct {
	width  int
	height int
}

type DrawLineCommand struct {
	x      int
	y      int
	length int
}

type DrawRectangleCommand struct {
	a int
	b int
	c int
	d int
}

type FillCommand struct {
	x           int
	y           int
	replacement rune
}

type QuitCommand struct{}

func ParseCommand(input string) (Command, error) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return nil, fmt.Errorf("el comando está vacío")
	}

	switch parts[0] {
	case "c":
		values, err := parseNumbers(parts[1:], 2)
		if err != nil {
			return nil, fmt.Errorf("uso: c width height: %w", err)
		}

		return CreateCanvasCommand{
			width:  values[0],
			height: values[1],
		}, nil

	case "l":
		values, err := parseNumbers(parts[1:], 3)
		if err != nil {
			return nil, fmt.Errorf("uso: l x y n: %w", err)
		}

		return DrawLineCommand{
			x:      values[0],
			y:      values[1],
			length: values[2],
		}, nil

	case "r":
		values, err := parseNumbers(parts[1:], 4)
		if err != nil {
			return nil, fmt.Errorf("uso: r a b c d: %w", err)
		}

		return DrawRectangleCommand{
			a: values[0],
			b: values[1],
			c: values[2],
			d: values[3],
		}, nil

	case "f":
		if len(parts) != 4 {
			return nil, fmt.Errorf("uso: f x y r")
		}

		values, err := parseNumbers(parts[1:3], 2)
		if err != nil {
			return nil, fmt.Errorf("uso: f x y r: %w", err)
		}

		replacement := []rune(parts[3])
		if len(replacement) != 1 {
			return nil, fmt.Errorf(
				"el reemplazo debe ser un solo carácter",
			)
		}

		return FillCommand{
			x:           values[0],
			y:           values[1],
			replacement: replacement[0],
		}, nil

	case "q":
		if len(parts) != 1 {
			return nil, fmt.Errorf("uso: q")
		}

		return QuitCommand{}, nil

	default:
		return nil, fmt.Errorf(
			"comando no reconocido: %s",
			parts[0],
		)
	}
}

func parseNumbers(
	values []string,
	expected int,
) ([]int, error) {
	if len(values) != expected {
		return nil, fmt.Errorf(
			"se esperaban %d argumentos y se recibieron %d",
			expected,
			len(values),
		)
	}

	numbers := make([]int, expected)

	for index, value := range values {
		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf(
				"%q no es un número válido: %w",
				value,
				err,
			)
		}

		numbers[index] = number
	}

	return numbers, nil
}

func (command CreateCanvasCommand) Execute(
	manager *CanvasManager,
) (CommandResult, error) {
	canvas := manager.Create(
		command.width,
		command.height,
	)

	return CommandResult{
		Canvas: canvas,
	}, nil
}

func (command DrawLineCommand) Execute(
	manager *CanvasManager,
) (CommandResult, error) {
	canvas, err := manager.Current()
	if err != nil {
		return CommandResult{}, err
	}

	err = canvas.DrawLine(
		command.x,
		command.y,
		command.length,
	)
	if err != nil {
		return CommandResult{}, fmt.Errorf(
			"no se pudo dibujar la línea: %w",
			err,
		)
	}

	return CommandResult{
		Canvas: canvas,
	}, nil
}

func (command DrawRectangleCommand) Execute(
	manager *CanvasManager,
) (CommandResult, error) {
	canvas, err := manager.Current()
	if err != nil {
		return CommandResult{}, err
	}

	err = canvas.DrawRectangle(
		command.a,
		command.b,
		command.c,
		command.d,
	)
	if err != nil {
		return CommandResult{}, fmt.Errorf(
			"no se pudo dibujar el rectángulo: %w",
			err,
		)
	}

	return CommandResult{
		Canvas: canvas,
	}, nil
}

func (command FillCommand) Execute(
	manager *CanvasManager,
) (CommandResult, error) {
	canvas, err := manager.Current()
	if err != nil {
		return CommandResult{}, err
	}

	err = canvas.Fill(
		command.x,
		command.y,
		command.replacement,
	)
	if err != nil {
		return CommandResult{}, fmt.Errorf(
			"no se pudo rellenar el canvas: %w",
			err,
		)
	}

	return CommandResult{
		Canvas: canvas,
	}, nil
}

func (QuitCommand) Execute(
	manager *CanvasManager,
) (CommandResult, error) {
	return CommandResult{
		Quit: true,
	}, nil
}
