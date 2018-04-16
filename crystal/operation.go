package crystal

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

var possiblePointSymmetrys [6960]PointOp

func init() {
	genAllPointOps()
}

func genAllPointOps() {
	vs := [26][]float64{
		[]float64{1, 0, 0},
		[]float64{0, 1, 0},
		[]float64{0, 0, 1},
		[]float64{-1, 0, 0},
		[]float64{0, -1, 0}, /* 5 */
		[]float64{0, 0, -1},
		[]float64{0, 1, 1},
		[]float64{1, 0, 1},
		[]float64{1, 1, 0},
		[]float64{0, -1, -1}, /* 10 */
		[]float64{-1, 0, -1},
		[]float64{-1, -1, 0},
		[]float64{0, 1, -1},
		[]float64{-1, 0, 1},
		[]float64{1, -1, 0}, /* 15 */
		[]float64{0, -1, 1},
		[]float64{1, 0, -1},
		[]float64{-1, 1, 0},
		[]float64{1, 1, 1},
		[]float64{-1, -1, -1}, /* 20 */
		[]float64{-1, 1, 1},
		[]float64{1, -1, 1},
		[]float64{1, 1, -1},
		[]float64{1, -1, -1},
		[]float64{-1, 1, -1}, /* 25 */
		[]float64{-1, -1, 1},
	}
	i := 0
	for _, up := range vs {
		for _, mid := range vs {
			for _, down := range vs {
				m := [][]float64{up, mid, down}
				if !isUniMat(m) {
					continue
				}
				possiblePointSymmetrys[i] = *NewPointOp(m)
				i++
			}
		}
	}
}

func isUniMat(m [][]float64) bool {
	det := m[0][0]*m[1][1]*m[2][2] +
		m[0][1]*m[1][2]*m[2][0] +
		m[0][2]*m[1][0]*m[2][1] -
		m[0][2]*m[1][1]*m[2][0] -
		m[0][1]*m[1][0]*m[2][2] -
		m[0][0]*m[1][2]*m[2][1]
	if det != 1.0 && det != -1.0 {
		return false
	}
	return true
}

type Rotator interface {
	Rotate(mat.Dense)
}

type PointOp struct {
	m mat.Dense
}

func NewPointOp(data [][]float64) *PointOp {
	shapeErr := fmt.Errorf("unexpect shape of rotation operation %v", data)
	if len(data) != 3 {
		panic(shapeErr)
	}
	for _, row := range data {
		if len(row) != 3 {
			panic(shapeErr)
		}
	}
	m := mat.NewDense(3, 3, nil)
	for i, src := range data {
		m.SetRow(i, src)
	}
	return &PointOp{
		*m,
	}
}

func (p PointOp) Det() float64 {
	return mat.Det(&p.m)
}

func (p PointOp) ActOn(r Rotator) {
	r.Rotate(p.m)
}

type Shifter interface {
	Shift(mat.VecDense)
}

type TransOp struct {
	v mat.VecDense
}

func NewTransOp(t []float64) *TransOp {
	v := mat.NewVecDense(3, t)
	return &TransOp{
		*v,
	}
}

func (t TransOp) ActOn(s Shifter) {
	s.Shift(t.v)
}

type Transformer interface {
	Rotator
	Shifter
}

type SpaceOp struct {
	R PointOp
	T TransOp
}

func (s SpaceOp) Transform(t Transformer) {
	t.Rotate(s.R.m)
	t.Shift(s.T.v)
}

type Extender interface {
	Extend(mat.Dense)
}

type HNF struct {
	m mat.TriDense
}

func NewHNF(n int, data []int) *HNF {
	fdata := make([]float64, len(data))
	for i := 0; i < len(fdata); i++ {
		fdata[i] = float64(data[i])
	}
	return &HNF{
		*mat.NewTriDense(n, false, fdata),
	}
}

func DupHnfs(v int) []HNF {
	r := make([]HNF, 0)
	acf := factor(v)
	for _, a := range acf {
		for _, c := range acf {
			f := v / a / c
			if f < 1 {
				continue
			}
			for b := 0; b < c; b++ {
				for d := 0; d < f; d++ {
					for e := 0; e < f; e++ {
						hm := *NewHNF(3, []int{a, 0, 0, b, c, 0, d, e, f})
						r = append(r, hm)
					}
				}
			}
		}
	}
	return r
}

func factor(nr int) []int {
	fs := make([]int, 1)
	fs[0] = 1
	apf := func(p int, e int) {
		n := len(fs)
		for i, pp := 0, p; i < e; i, pp = i+1, pp*p {
			for j := 0; j < n; j++ {
				fs = append(fs, fs[j]*pp)
			}
		}
	}
	e := 0
	for ; nr&1 == 0; e++ {
		nr >>= 1
	}
	apf(2, e)
	for d := int(3); nr > 1; d += 2 {
		if d*d > nr {
			d = nr
		}
		for e = 0; nr%d == 0; e++ {
			nr /= d
		}
		if e > 0 {
			apf(d, e)
		}
	}
	return fs
}
