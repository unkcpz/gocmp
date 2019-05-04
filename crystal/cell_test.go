package crystal

import (
  // "fmt"
  "testing"
)

func TestNewCell(t *testing.T) {
  _, err := NewCell(
    []float64{0,0,0,0,0,0,0,0,0},
    []float64{0,0,0,0,0,0},
    []int{0,0},
  )
  if err == nil {
    t.Errorf("expect NewCell error")
  }
}
