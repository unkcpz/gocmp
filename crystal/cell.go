package crystal

import (
	"fmt"
	"math"

	"github.com/unkcpz/spgolib"
	"gonum.org/v1/gonum/mat"
)

type Cell struct {
	// Lattice row vectors
	Lattice *mat.Dense
	// Position Fraction
	Position *mat.Dense
	// atom types of cell
	Elem []int
	// number of atoms
	Natom int
}

func CellCopyOf(c *Cell) *Cell {
  latt := mat.DenseCopyOf(c.Lattice)
  pos := mat.DenseCopyOf(c.Position)
  n := c.Natom
  elem := make([]int, n, n)
  for i:=0; i<n; i++ {
    elem[i] = c.Elem[i]
  }
	r := &Cell{
		Lattice:  latt,
		Position: pos,
		Elem:     elem,
		Natom:    n,
	}
  return r
}

// NewCell create cell
func NewCell(
	lattice []float64,
	position []float64,
	types []int,
	cartesian bool) (*Cell, error) {

	if len(lattice) != 9 {
		return nil, fmt.Errorf("expext 9 value as lattice")
	}
	latt := mat.NewDense(3, 3, lattice)
	if math.Abs(mat.Det(latt)) < 1e-5 {
		return nil, fmt.Errorf("lattice cannot be det=0")
	}

	if len(position) != 3*len(types) {
		return nil, fmt.Errorf("atom number not compatible, len(types)=%d, len(positions)/3=%d/3", len(types), len(position))
	}
	n := len(types)
	pos := mat.NewDense(n, 3, position)
	if cartesian {
		var ia mat.Dense
		ia.Inverse(latt)
		pos.Mul(pos, &ia)
	}

	cell := &Cell{
		Lattice:  latt,
		Position: pos,
		Elem:     types,
		Natom:    n,
	}
	return cell, nil
}

// LatticeSlice return lattice as []float64
func (c *Cell) LatticeSlice() []float64 {
	return matToSlice(c.Lattice)
}

// PositionSlice return position as []float64
func (c *Cell) PositionSlice() []float64 {
	return matToSlice(c.Position)
}

// Primitive return primitive cell
func (c *Cell) Primitive(symprec float64) {
	lattice := c.LatticeSlice()
	position := c.PositionSlice()
	elem := c.Elem
	n := c.Natom
	newLattice, newPosition, newElem :=
		spgolib.Standardize(lattice, position, elem, n, true, true, symprec)
	c.Lattice = mat.NewDense(3, 3, newLattice)
	c.Position = mat.NewDense(len(newElem), 3, newPosition)
	c.Elem = newElem
	c.Natom = len(newElem)
}

// Refine return conventional cell
func (c *Cell) Refine(symprec float64) {
	lattice := c.LatticeSlice()
	position := c.PositionSlice()
	elem := c.Elem
	n := c.Natom
	newLattice, newPosition, newElem :=
		spgolib.Standardize(lattice, position, elem, n, false, false, symprec)
	c.Lattice = mat.NewDense(3, 3, newLattice)
	c.Position = mat.NewDense(len(newElem), 3, newPosition)
	c.Elem = newElem
	c.Natom = len(newElem)
}

// Symmetry find rotation and translation of cell
func (c *Cell) Symmetry(symprec float64) (nop int, rotations []Rotation, transitions []Translation) {
  ds := spgolib.NewDataset(c.LatticeSlice(), c.PositionSlice(), c.Elem, symprec)
  nop = ds.Nops
  rots := make([]Rotation, nop, nop)
  trans := make([]Translation, nop, nop)
  for i:=0; i<nop; i++ {
    rots[i].data = mat.NewDense(3, 3, intToFloat(ds.Rotations[i*9:i*9+9]))
    trans[i].data = mat.NewVecDense(3, ds.Translations[i*3:i*3+3])
  }
  return nop, rots, trans
}

func (c *Cell) String() string {
  fa := mat.Formatted(c.Lattice, mat.Squeeze())
  fb := mat.Formatted(c.Position, mat.Squeeze())
  return fmt.Sprintf("%v\n%v", fa, fb)
}

func intToFloat(a []int) []float64 {
	b := make([]float64, len(a))
	for i := range b {
		b[i] = float64(a[i])
	}
	return b
}

func matToSlice(mat *mat.Dense) []float64 {
	blasM := mat.RawMatrix()
	return blasM.Data
}
