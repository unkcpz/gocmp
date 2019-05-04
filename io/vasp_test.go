package io

import(
  "testing"
)

func TestPoscarRead(t *testing.T) {
  txt := `system
1
4.0 0.0 0.0
0.0 4.0 0.0
0.0 0.0 4.0
H
2
selective dynamics
Karti
0.0 0.0 0.0
2.0 2.0 2.0`
  poscar, err := ParsePoscar(txt)
  if err != nil {
    t.Error("Parse error")
  }

  expect_latt := []float64{4.0, 0.0, 0.0, 0.0, 4.0, 0.0, 0.0, 0.0, 4.0}
  for i, _ := range expect_latt {
    if poscar.Lattice[i] - expect_latt[i] > 1e-5 {
      t.Error("poscar lattice read fail")
      break
    }
  }

  expect_pos := []float64{0.0, 0.0, 0.0, 2.0, 2.0, 2.0}
  for i, _ := range expect_pos {
    if poscar.Positions[i] - expect_pos[i] > 1e-5 {
      t.Error("poscar positions read fail")
      break
    }
  }

  expect_types := []string{"H", "H"}
  for i, _ := range expect_types {
    if poscar.Types[i] != expect_types[i] {
      t.Error("poscar types read fail")
      break
    }
  }

  if poscar.Coordinate != Cartesian {
    t.Error("poscar coordinate type parse failed")
  }
}
