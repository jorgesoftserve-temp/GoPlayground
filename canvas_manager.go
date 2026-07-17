package main

import "errors"

type CanvasManager struct {
	current *Canvas
}

func NewCanvasManager() *CanvasManager {
	return &CanvasManager{}
}

func (manager *CanvasManager) Create(
	width int,
	height int,
) *Canvas {
	manager.current = NewCanvas(width, height)

	return manager.current
}

func (manager *CanvasManager) Current() (*Canvas, error) {
	if manager.current == nil {
		return nil, errors.New(
			"primero crea un canvas con: c width height",
		)
	}

	return manager.current, nil
}
