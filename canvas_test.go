package main

import "testing"

func TestNewCanvas(t *testing.T) {
	canvas := NewCanvas(3, 2)

	expected := "" +
		"-----\n" +
		"|   |\n" +
		"|   |\n" +
		"-----\n"

	result := canvas.String()

	if result != expected {
		t.Fatalf(
			"canvas inesperado\nesperado:\n%s\nobtenido:\n%s",
			expected,
			result,
		)
	}
}

func TestDrawLine(t *testing.T) {
	canvas := NewCanvas(5, 3)

	err := canvas.DrawLine(2, 2, 3)
	if err != nil {
		t.Fatalf("DrawLine devolvió error: %v", err)
	}

	expected := "" +
		"-------\n" +
		"|     |\n" +
		"| XXX |\n" +
		"|     |\n" +
		"-------\n"

	if result := canvas.String(); result != expected {
		t.Fatalf(
			"línea inesperada\nesperado:\n%s\nobtenido:\n%s",
			expected,
			result,
		)
	}
}

func TestDrawRectangle(t *testing.T) {
	canvas := NewCanvas(6, 4)

	err := canvas.DrawRectangle(2, 2, 5, 4)
	if err != nil {
		t.Fatalf("DrawRectangle devolvió error: %v", err)
	}

	expected := "" +
		"--------\n" +
		"|      |\n" +
		"| XXXX |\n" +
		"| X  X |\n" +
		"| XXXX |\n" +
		"--------\n"

	if result := canvas.String(); result != expected {
		t.Fatalf(
			"rectángulo inesperado\nesperado:\n%s\nobtenido:\n%s",
			expected,
			result,
		)
	}
}

func TestFillOnlyChangesConnectedArea(t *testing.T) {
	canvas := NewCanvas(6, 4)

	err := canvas.DrawRectangle(2, 2, 5, 4)
	if err != nil {
		t.Fatalf("DrawRectangle devolvió error: %v", err)
	}

	err = canvas.Fill(1, 1, '.')
	if err != nil {
		t.Fatalf("Fill devolvió error: %v", err)
	}

	expected := "" +
		"--------\n" +
		"|......|\n" +
		"|.XXXX.|\n" +
		"|.X  X.|\n" +
		"|.XXXX.|\n" +
		"--------\n"

	if result := canvas.String(); result != expected {
		t.Fatalf(
			"fill inesperado\nesperado:\n%s\nobtenido:\n%s",
			expected,
			result,
		)
	}
}

func TestFillInsideRectangle(t *testing.T) {
	canvas := NewCanvas(6, 5)

	err := canvas.DrawRectangle(2, 2, 5, 4)
	if err != nil {
		t.Fatalf("DrawRectangle devolvió error: %v", err)
	}

	err = canvas.Fill(3, 3, 'o')
	if err != nil {
		t.Fatalf("Fill devolvió error: %v", err)
	}

	expected := "" +
		"--------\n" +
		"|      |\n" +
		"| XXXX |\n" +
		"| XooX |\n" +
		"| XXXX |\n" +
		"|      |\n" +
		"--------\n"

	if result := canvas.String(); result != expected {
		t.Fatalf(
			"fill interior inesperado\nesperado:\n%s\nobtenido:\n%s",
			expected,
			result,
		)
	}
}

func TestGetPixelOutsideCanvasReturnsError(t *testing.T) {
	canvas := NewCanvas(3, 3)

	_, err := canvas.GetPixel(5, 5)
	if err == nil {
		t.Fatal("se esperaba un error para una coordenada inválida")
	}
}
