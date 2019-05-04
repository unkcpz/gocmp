package io

import (
  "strings"
  "fmt"
  "strconv"
)

type CoorType int

const (
  Cartesian CoorType = iota
  Fractional
)

type Poscar struct {
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

func ParsePoscar(txt string) (*Poscar, error) {
  lines := strings.Split(txt, "\n")

  system := lines[0]

  scale, err := strconv.ParseFloat(lines[1], 64)
  if err != nil {
    return nil, fmt.Errorf("parse poscar scale, %v", err)
  }

  lattice := make([]float64, 9, 9)
  for i:=0; i<3; i++ {
    vs := strings.Split(lines[2+i], " ")
    lattice[i*3+0], _ = strconv.ParseFloat(vs[0], 64)
    lattice[i*3+1], _ = strconv.ParseFloat(vs[1], 64)
    lattice[i*3+2], _ = strconv.ParseFloat(vs[2], 64)
  }
  for i, _ := range lattice {
    lattice[i] *= scale
  }

  es := strings.Split(lines[5], " ")
  ntype := make([]int, 0, 0)
  for _, s := range strings.Split(lines[6], " ") {
    n, _ := strconv.Atoi(s)
    ntype = append(ntype, n)
  }
  var natoms int
  for _, v := range ntype {
    natoms += v
  }
  types := make([]string, 0, 0)
  for i, n := range ntype {
    for j:=0; j<n; j++ {
      types = append(types, es[i])
    }
  }

  var coorLine int = 7
  s := strings.ToUpper(lines[7][:1])
  if s == "S" {
    coorLine = 8
  }

  var ctype CoorType = Fractional
  ct := strings.ToUpper(lines[coorLine][:1])
  if ct == "K" || ct == "C" {
    ctype = Cartesian
  }

  positions := make([]float64, 3*natoms, 3*natoms)
  for i:=0; i<natoms; i++ {
    vs := strings.Split(lines[coorLine+1+i], " ")
    positions[i*3+0], _ = strconv.ParseFloat(vs[0], 64)
    positions[i*3+1], _ = strconv.ParseFloat(vs[1], 64)
    positions[i*3+2], _ = strconv.ParseFloat(vs[2], 64)
  }

  poscar := &Poscar {
    System: system,
    Lattice: lattice,
    Coordinate: ctype,
    Positions: positions,
    Types: types,
  }
  return poscar, nil
}
