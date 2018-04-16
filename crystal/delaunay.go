package crystal

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// DealaunayReduce reduce lattice to standard lattice
func DealaunayReduce(latt *Lattice, prec float64) (*Lattice, error) {
	bs := make(BasisSlice, 7)
	bs[0], bs[1], bs[2] = latt.A(), latt.B(), latt.C()
	bs[3] = bs[0].Add(bs[1]).Add(bs[2]).Scale(-1) // b3 = -b0-b1-b2

	// reduce loop
	for i := 0; i < 4; i++ {
		for j := 1; j < 4; j++ {
			if bs[i].Dot(bs[j]) < prec {
				break
			}
			for e := 0; e < 4; e++ {
				if e != i && e != j {
					bs[e] = bs[e].Add(bs[i])
				}
			}
			bs[i] = bs[i].Scale(-1)
		}
	}

	// shortest vector
	bs[4] = bs[0].Add(bs[1])
	bs[5] = bs[1].Add(bs[2])
	bs[6] = bs[2].Add(bs[0])
	sort.Sort(bs)

	// l := NewLatticeR(*bs[0], *bs[1], *NewBasis(nil))
	m := mat.NewDense(3, 3, nil)
	m.SetRow(0, bs[0].Data())
	m.SetRow(1, bs[1].Data())
	for i := 2; i < 7; i++ {
		m.SetRow(2, bs[i].Data())
		if math.Abs(mat.Det(m)) > prec {
			break
		}
	}

	// right hand axis
	if mat.Det(m) < -prec {
		m.Scale(-1, m)
	}

	return &Lattice{
		*m,
	}, nil
}
