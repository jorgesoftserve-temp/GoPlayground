package main

import (
	"fmt"
	"strings"
)

type Canvas struct {
	width  int
	height int
	pixels [][]rune
}

type Point struct {
	x int
	y int
}

func NewCanvas(width int, height int) *Canvas {
	pixels := make([][]rune, height)

	for y := range height {
		pixels[y] = make([]rune, width)

		for x := range width {
			pixels[y][x] = ' '
		}
	}

	return &Canvas{
		width:  width,
		height: height,
		pixels: pixels,
	}
}

func (canvas *Canvas) String() string {
	var builder strings.Builder

	for range canvas.width + 2 {
		builder.WriteRune('-')
	}
	builder.WriteRune('\n')

	for y := range canvas.pixels {
		builder.WriteRune('|')

		for x := range canvas.pixels[y] {
			builder.WriteRune(canvas.pixels[y][x])
		}

		builder.WriteString("|\n")
	}

	for range canvas.width + 2 {
		builder.WriteRune('-')
	}
	builder.WriteRune('\n')

	return builder.String()
}

func (canvas *Canvas) DrawLine(x int, y int, length int) error {
	startX := x - 1
	startY := y - 1

	for i := 0; i < length; i++ {
		err := canvas.SetPixel(startX+i, startY, 'X')
		if err != nil {
			return err
		}
	}

	return nil
}

func (canvas *Canvas) DrawRectangle(
	a int,
	b int,
	c int,
	d int,
) error {
	left := a - 1
	top := b - 1
	right := c - 1
	bottom := d - 1

	for x := left; x <= right; x++ {
		if err := canvas.SetPixel(x, top, 'X'); err != nil {
			return err
		}

		if err := canvas.SetPixel(x, bottom, 'X'); err != nil {
			return err
		}
	}

	for y := top; y <= bottom; y++ {
		if err := canvas.SetPixel(left, y, 'X'); err != nil {
			return err
		}

		if err := canvas.SetPixel(right, y, 'X'); err != nil {
			return err
		}
	}

	return nil
}

func (canvas *Canvas) Fill(
	x int,
	y int,
	replacement rune,
) error {
	start := Point{
		x: x - 1,
		y: y - 1,
	}

	target, err := canvas.GetPixel(start.x, start.y)
	if err != nil {
		return err
	}

	if target == replacement {
		return nil
	}

	points := []Point{start}
	visited := map[Point]bool{
		start: true,
	}

	for len(points) > 0 {
		lastIndex := len(points) - 1
		current := points[lastIndex]
		points = points[:lastIndex]

		currentValue, err := canvas.GetPixel(current.x, current.y)
		if err != nil {
			continue
		}

		if currentValue != target {
			continue
		}

		if err := canvas.SetPixel(
			current.x,
			current.y,
			replacement,
		); err != nil {
			return err
		}

		neighbors := []Point{
			{x: current.x + 1, y: current.y},
			{x: current.x - 1, y: current.y},
			{x: current.x, y: current.y + 1},
			{x: current.x, y: current.y - 1},
		}

		for _, neighbor := range neighbors {
			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			points = append(points, neighbor)
		}
	}

	return nil
}

func (canvas *Canvas) GetPixel(x int, y int) (rune, error) {
	if x < 0 || x >= canvas.width {
		return 0, fmt.Errorf(
			"coordenada x fuera del canvas: %d",
			x+1,
		)
	}

	if y < 0 || y >= canvas.height {
		return 0, fmt.Errorf(
			"coordenada y fuera del canvas: %d",
			y+1,
		)
	}

	return canvas.pixels[y][x], nil
}

func (canvas *Canvas) SetPixel(
	x int,
	y int,
	char rune,
) error {
	if x < 0 || x >= canvas.width {
		return fmt.Errorf(
			"coordenada x fuera del canvas: %d",
			x+1,
		)
	}

	if y < 0 || y >= canvas.height {
		return fmt.Errorf(
			"coordenada y fuera del canvas: %d",
			y+1,
		)
	}

	canvas.pixels[y][x] = char

	return nil
}
