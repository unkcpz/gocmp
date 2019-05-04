package crystal

import (
  // "fmt"
  "math"
  "testing"
  
  "gonum.org/v1/gonum/mat"
)

func TestNewCell(t *testing.T) {
  _, err := NewCell(
    []float64{0,0,0,0,0,0,0,0,0},
    []float64{0,0,0,0,0,0},
    []int{0,0},
    false,
  )
  if err == nil {
    t.Errorf("expect NewCell error")
  }

  c, _ := NewCell(
    []float64{4,0,0,0,4,0,0,0,4},
    []float64{0,0,0,2,2,2},
    []int{1, 1},
    true,
  )
  expect_pos := mat.NewDense(2, 3, []float64{0,0,0,0.5,0.5,0.5})
  for i:=0; i<2; i++ {
    for j:=0; j<3; j++ {
      if math.Abs(c.Positions.At(i, j) - expect_pos.At(i, j)) > 1e-5 {
        t.Error("NewCell for cartesian error")
        break
      }
    }
  }
}
