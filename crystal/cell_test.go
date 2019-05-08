package crystal

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNewCell(t *testing.T) {
	_, err := NewCell(
		[]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float64{0, 0, 0, 0, 0, 0},
		[]int{0, 0},
		false,
	)
	if err == nil {
		t.Errorf("expect NewCell error")
	}

	// test init with cartesian coordinate
	c, _ := NewCell(
		[]float64{4, 0, 0, 0, 4, 0, 0, 0, 4},
		[]float64{0, 0, 0, 2, 2, 2},
		[]int{1, 1},
		true,
	)
	expect_pos := mat.NewDense(2, 3, []float64{0, 0, 0, 0.5, 0.5, 0.5})
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(c.Position.At(i, j)-expect_pos.At(i, j)) > 1e-5 {
				t.Error("NewCell for cartesian error")
				break
			}
		}
	}
}

func TestCellCopyOf(t *testing.T) {
	c, _ := NewCell(
		[]float64{4, 0, 0, 0, 4, 0, 0, 0, 4},
		[]float64{0, 0, 0, 2, 2, 2},
		[]int{1, 1},
		true,
	)
  newc := CellCopyOf(c)

  if &newc == &c || &newc.Lattice == &c.Lattice ||
      &newc.Position == &c.Position || &newc.Elem == &c.Elem {
    t.Error("expected different object")
  }
}

func TestPrimitive(t *testing.T) {
	c, _ := NewCell(
		[]float64{4, 0, 0, 0, 4, 0, 0, 0, 4},
		[]float64{0, 0, 0, 0.5, 0.5, 0.5},
		[]int{1, 1},
		false,
	)

  c.Primitive(1e-5)
	expected, _ := NewCell(
		[]float64{-2, 2, 2, 2, -2, 2, 2, 2, -2},
		[]float64{0, 0, 0},
		[]int{1},
		false,
	)
	if !mat.EqualApprox(c.Lattice, expected.Lattice, 1e-5) {
		t.Errorf("lattice: primitive expected %v, got %v.", expected, c)
	}
	if !mat.EqualApprox(c.Position, expected.Position, 1e-5) {
		t.Errorf("position: primitive expected %v, got %v.", expected, c)
	}
  for i, _ := range c.Elem {
  	if c.Elem[i] != expected.Elem[i] {
  		t.Errorf("elem: primitive expected %v, got %v.", expected, c)
  	}
  }
}

func TestRefine(t *testing.T) {
	c, _ := NewCell(
		[]float64{-2, 2, 2, 2, -2, 2, 2, 2, -2},
		[]float64{0, 0, 0},
		[]int{1},
		false,
	)
  c.Refine(1e-5)

	expected, _ := NewCell(
		[]float64{4, 0, 0, 0, 4, 0, 0, 0, 4},
		[]float64{0, 0, 0, 0.5, 0.5, 0.5},
		[]int{1, 1},
		false,
	)

	if !mat.EqualApprox(c.Lattice, expected.Lattice, 1e-5) {
		t.Errorf("lattice: primitive expected %v, got %v.", expected, c)
	}
	if !mat.EqualApprox(c.Position, expected.Position, 1e-5) {
		t.Errorf("position: primitive expected %v, got %v.", expected, c)
	}
  for i, _ := range c.Elem {
  	if c.Elem[i] != expected.Elem[i] {
  		t.Errorf("elem: primitive expected %v, got %v.", expected, c)
  	}
  }
}
