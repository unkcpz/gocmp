package crystal

import (
  "fmt"
  "math"

  "gonum.org/v1/gonum/mat"
)

type Cell struct {
  // Lattice
  Lattice *mat.Dense
  // Position Cartesian
  Positions *mat.Dense
  //
  Types []int
  //
  Natoms int
}

func NewCell(
  lattice []float64,
  positions []float64,
  types []int,
  cartesian bool) (*Cell, error) {

  if len(lattice) != 9 {
    return nil, fmt.Errorf("expext 9 value as lattice")
  }
  latt := mat.NewDense(3, 3, lattice)
  if math.Abs(mat.Det(latt)) < 1e-5 {
    return nil, fmt.Errorf("lattice cannot be det=0")
  }

  if len(positions) != 3*len(types) {
    return nil, fmt.Errorf("atom number not compatible, len(types)=%d, len(positions)/3=%d/3", len(types), len(positions))
  }
  n := len(types)
  pos := mat.NewDense(n, 3, positions)
  if cartesian {
    var ia mat.Dense
    ia.Inverse(latt)
    pos.Mul(pos, &ia)
  }

  cell := &Cell {
    Lattice: latt,
    Positions: pos,
    Types: types,
    Natoms: n,
  }
  return cell, nil
}
