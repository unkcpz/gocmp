// crystal package contains fundamental data structures
// and fundamental operations for crystals
package crystal

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// Basis type represent crystal basis in Lattice
type Basis struct {
	v mat.VecDense
}

// NewBasis create a new Basis
func NewBasis(v []float64) *Basis {
	r := mat.NewVecDense(3, v)
	return &Basis{
		*r,
	}
}

func (b *Basis) Data() []float64 {
	return b.v.RawVector().Data
}

// Add plus bx into b
func (b *Basis) Add(bx *Basis) *Basis {
	r := mat.NewVecDense(3, nil)
	r.AddVec(&b.v, &bx.v)
	return &Basis{
		*r,
	}
}

// Sub minux bx from b
func (b *Basis) Sub(bx *Basis) *Basis {
	r := mat.NewVecDense(3, nil)
	r.SubVec(&b.v, &bx.v)
	return &Basis{
		*r,
	}
}

// Mul return the vector point multiple result
func (b *Basis) Dot(bx *Basis) float64 {
	return mat.Dot(&b.v, &bx.v)
}

// Scale multi each value in basis
func (b *Basis) Scale(s float64) *Basis {
	r := mat.NewVecDense(3, nil)
	r.ScaleVec(s, &b.v)
	return &Basis{
		*r,
	}
}

// Norm get the square norm of basis
func (b *Basis) Norm() float64 {
	return mat.Norm(&b.v, 2)
}

type BasisSlice []*Basis

func (x BasisSlice) Len() int {
	return len(x)
}

func (x BasisSlice) Less(i, j int) bool {
	return x[i].Norm() < x[j].Norm()
}

func (x BasisSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type Lattice struct {
	m mat.Dense
}

// NewLattice create a new Lattice
func NewLattice(l []float64) *Lattice {
	return &Lattice{
		*mat.NewDense(3, 3, l),
	}
}

func (l *Lattice) RawData() []float64 {
	return l.m.RawMatrix().Data
}

func (l *Lattice) DeepCopy() *Lattice {
	return &Lattice{
		*mat.DenseCopyOf(&l.m),
	}
}

func (l *Lattice) A() *Basis {
	return NewBasis(l.m.RawRowView(0))
}

func (l *Lattice) B() *Basis {
	return NewBasis(l.m.RawRowView(1))
}

func (l *Lattice) C() *Basis {
	return NewBasis(l.m.RawRowView(2))
}

// Det get the detereminant of the lattice
func (l *Lattice) Det() float64 {
	return mat.Det(&l.m)
}

// T get the transpose of Lattice
func (l *Lattice) T() *Lattice {
	m := mat.NewDense(3, 3, nil)
	for i := 0; i < 3; i++ {
		m.SetCol(i, l.m.RawRowView(i))
	}
	return &Lattice{
		*m,
	}
}

func (l *Lattice) Rotate(m mat.Dense) {
	l.m.Mul(&m, &l.m)
}

// Gram get the gram matrix of lattice
func (l *Lattice) Gram() *mat.Dense {
	m := mat.NewDense(3, 3, nil)
	m.Mul(&l.m, l.m.T())
	return m
}

// Isometry find is two lattice congrent
func (l *Lattice) Equal(lx *Lattice, prec float64) bool {
	return mat.EqualApprox(l.Gram(), lx.Gram(), prec)
}

// Symmetry returns symmetrys of lattice
func (l *Lattice) Symmetry(prec float64) []PointOp {
	var ps []PointOp
	// dl, err := DealaunayReduce(l, prec)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(*dl)

	for _, r := range possiblePointSymmetrys {
		newl := l.DeepCopy()
		newl.Rotate(r.m)
		if newl.Equal(l, math.Sqrt(prec)) {
			ps = append(ps, r)
		}
	}
	return ps
}

type FPos struct {
	v mat.VecDense
}

func NewFPos(data []float64) *FPos {
	r := mat.NewVecDense(3, data)
	return &FPos{
		*r,
	}
}

func (fp *FPos) RawData() []float64 {
	return fp.v.RawVector().Data
}

func (fp *FPos) Sub(v mat.Vector) *FPos {
	r := mat.NewVecDense(3, nil)
	r.SubVec(&fp.v, v)
	return &FPos{
		*r,
	}
}

func (fp *FPos) Add(v mat.Vector) *FPos {
	r := mat.NewVecDense(3, nil)
	r.AddVec(&fp.v, v)
	return &FPos{
		*r,
	}
}

func (fp *FPos) EqualInCell(fpx *FPos, prec float64) bool {
	return fp.InCell(prec).Equal(fpx.InCell(prec), prec)
}

func (fp *FPos) Equal(fpx *FPos, prec float64) bool {
	return mat.EqualApprox(&fp.v, &fpx.v, prec)
}

func (fp *FPos) InCell(prec float64) *FPos {
	r := mat.NewVecDense(3, nil)
	for i := 0; i < fp.v.Len(); i++ {
		val := func(v float64) float64 {
			v = math.Mod(v, 1)
			if v < -prec {
				v += 1
			}
			if v < 0 && v > -prec {
				v = 0
			}
			return v
		}(fp.v.AtVec(i))
		r.SetVec(i, val)
	}
	return &FPos{
		*r,
	}
}

type EType []string

type Cell struct {
	Lattice   Lattice
	Fractions mat.Dense
	EType     EType

	Len int
}

func NewCell(l []float64, f []float64, et []string) *Cell {
	n := len(et)
	return &Cell{
		*NewLattice(l),
		*mat.NewDense(n, 3, f),
		et,
		n,
	}
}

func (c *Cell) DeepCopy() *Cell {
	// can't just use RawData of c
	// because []float64 use as a refrence
	// newcf := c.Fractions.DeepCopy()
	copylm := mat.DenseCopyOf(&c.Lattice.m)
	copyfm := mat.DenseCopyOf(&c.Fractions)
	copyet := make(EType, len(c.EType))
	copy(copyet, c.EType)
	return NewCell(copylm.RawMatrix().Data, copyfm.RawMatrix().Data, copyet)
}

func (c *Cell) Idx(i int) *FPos {
	return NewFPos(c.Fractions.RawRowView(i))
}

func (c *Cell) Symmetry(prec float64) []SpaceOp {
	var sops []SpaceOp
	lti := c.IdxOfEle(c.leastEle())

	pOps := c.Lattice.Symmetry(prec)
	for _, r := range pOps {
		newc := c.DeepCopy()
		r.ActOn(newc)
		trs := func(co *Cell, cx *Cell, lti []int) []TransOp {
			trs := make([]TransOp, 0)
			o := co.Idx(lti[0])
			for _, i := range lti {
				tv := o.Sub(&cx.Idx(i).v).InCell(prec)
				t := NewTransOp(tv.RawData())
				trs = append(trs, *t)
			}
			return trs
		}(c, newc, lti)
		for _, t := range trs {
			nc := newc.DeepCopy()
			t.ActOn(nc)
			if ok := nc.Equal(c, prec); !ok {
				continue
			}
			var sop SpaceOp
			sop.R = r
			sop.T = t
			sops = append(sops, sop)
		}
	}
	return sops
}

func (c *Cell) Equal(cx *Cell, prec float64) bool {
	if len(c.EType) != len(c.EType) {
		return false
	}
	if !c.Lattice.Equal(&cx.Lattice, prec) {
		return false
	}
	r := c.Len
	for i := 0; i < r; i++ {
		v, e := c.Idx(i), c.EType[i]
		if _, ok := cx.Find(v, e, prec); !ok {
			return false
		}
	}
	return true
}

func (c *Cell) Find(v *FPos, e string, prec float64) (idx int, ok bool) {
	for i := 0; i < len(c.EType); i++ {
		if c.EType[i] == e && c.Idx(i).EqualInCell(v, prec) {
			return i, true
		}
	}
	return 0, false
}

func (c *Cell) Rotate(m mat.Dense) {
	c.Fractions.Mul(&c.Fractions, &m)
}

func (c *Cell) Shift(v mat.VecDense) {
	// r, _ := fps.Dims()
	r := c.Len
	for i := 0; i < r; i++ {
		vd := mat.NewVecDense(3, nil)
		vd.AddVec(c.Fractions.RowView(i), &v)
		c.Fractions.SetRow(i, vd.RawVector().Data)
	}
	// c.Fractions.Shift(v)
}

func (c *Cell) Extend(m mat.Dense) {

}

func (c *Cell) IdxOfEle(e string) []int {
	var r []int
	for i, t := range c.EType {
		if t == e {
			r = append(r, i)
		}
	}
	return r
}

func (c *Cell) leastEle() string {
	mes := make(map[string]int)
	for _, e := range c.EType {
		mes[e]++
	}

	var keys []string
	for k := range mes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var ele = keys[0]
	var min = mes[keys[0]]
	for i := 1; i < len(keys); i++ {
		if m := mes[keys[i]]; m < min {
			min = m
			ele = keys[i]
		}
	}
	return ele
}
