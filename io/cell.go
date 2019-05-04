package io

type CoorType int

const (
  Cartesian CoorType = iota
  Fractional
)

type Cell struct {
  //
  System string
  //
  Lattice []float64
  //
  Coordinate CoorType
  //
  Positions []float64
  //
  Types []string
}
