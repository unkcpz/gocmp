package crystal

import (
  "fmt"
)

func ExampleGetSpacegroup() {
  // Triangle lattice
	c, _ := NewCell(
		[]float64{2, 0, 0, 1, 1.732050, 0, 0, 0, 16},
		[]float64{0, 0, 0},
		[]int{1},
		false,
	)
  fmt.Println(c.Spacegroup(1e-5))

  // Output:
  // P6/mmm (191)
}

func ExampleGetSymmetry() {
  // Triangle lattice
	c, _ := NewCell(
		[]float64{2, 0, 0, 1, 1.732050, 0, 0, 0, 16},
		[]float64{0, 0, 0},
		[]int{1},
		false,
	)
  n, rots, trans := c.Symmetry(1e-5)
  fmt.Println(n)
  for _, r := range rots[:4] {
    fmt.Println(r)
  }
  for _, tr := range trans[:4] {
    fmt.Println(tr)
  }

  // Output:
  // 24
  // ⎡1  0  0⎤
  // ⎢0  1  0⎥
  // ⎣0  0  1⎦
  // ⎡-1   0   0⎤
  // ⎢ 0  -1   0⎥
  // ⎣ 0   0  -1⎦
  // ⎡0  -1  0⎤
  // ⎢1   1  0⎥
  // ⎣0   0  1⎦
  // ⎡ 0   1   0⎤
  // ⎢-1  -1   0⎥
  // ⎣ 0   0  -1⎦
  // [0  -0  0]
  // [0  -0  0]
  // [0  -0  0]
  // [0  -0  0]
}

func ExampleCell() {
  // Triangle lattice
	c, _ := NewCell(
		[]float64{2, 0, 0, 1, 1.732050, 0, 0, 0, 16},
		[]float64{0, 0, 0},
		[]int{1},
		false,
	)
  fmt.Println(c)

  // Output:
  // ⎡2        0   0⎤
  // ⎢1  1.73205   0⎥
  // ⎣0        0  16⎦
  // [0  0  0]

}
