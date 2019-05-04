package crystal

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNewBasis(t *testing.T) {
	for i, test := range []struct {
		input  []float64
		wanted *Basis
	}{
		{input: []float64{1.2, 0.8, 0.2}, wanted: NewBasis([]float64{1.2, 0.8, 0.2})},
		{input: nil, wanted: NewBasis([]float64{0.0, 0.0, 0.0})},
	} {
		b := NewBasis(test.input)
		if got := b; !reflect.DeepEqual(got, test.wanted) {
			t.Errorf("unexpected Basis data for test %d, got %v, want %v", i, got, test.wanted)
		}
	}
}

func TestBasisAdd(t *testing.T) {
	var tests = []struct {
		b      *Basis
		bx     *Basis
		wanted *Basis
	}{
		{
			b:      NewBasis([]float64{1, 3, 5}),
			bx:     NewBasis([]float64{2, 0, -2}),
			wanted: NewBasis([]float64{3, 3, 3.0}),
		},
	}

	for i, test := range tests {
		if got := test.b.Add(test.bx); !reflect.DeepEqual(got, test.wanted) {
			t.Errorf("test %d: %v.Add(%v) = %v, want %v", i, test.b, test.bx, got, test.wanted)
		}
	}
}

func TestBasisSub(t *testing.T) {
	var tests = []struct {
		b      *Basis
		bx     *Basis
		wanted *Basis
	}{
		{
			b:      NewBasis([]float64{1, 3, 5}),
			bx:     NewBasis([]float64{2, 0, -2}),
			wanted: NewBasis([]float64{-1, 3, 7.0}),
		},
	}

	for i, test := range tests {
		if got := test.b.Sub(test.bx); !reflect.DeepEqual(got, test.wanted) {
			t.Errorf("test %d: %v.Sub(%v) = %v, want %v", i, test.b, test.bx, got, test.wanted)
		}
	}
}

func TestBasisDot(t *testing.T) {
	var tests = []struct {
		b      *Basis
		bx     *Basis
		wanted float64
	}{
		{
			b:      NewBasis([]float64{1, 3, 5}),
			bx:     NewBasis([]float64{2, 0, -2}),
			wanted: -8.0,
		},
	}

	for i, test := range tests {
		if got := test.b.Dot(test.bx); got != test.wanted {
			t.Errorf("test %d: %v.Mul(%v) = %f, want %f", i, test.b, test.bx, got, test.wanted)
		}
	}
}

func TestBasisScale(t *testing.T) {
	var tests = []struct {
		b      *Basis
		scale  float64
		wanted *Basis
	}{
		{
			b:      NewBasis([]float64{1.1, 0, -2.1}),
			scale:  -1.0,
			wanted: NewBasis([]float64{-1.1, 0, 2.1}),
		},
		{
			b:      NewBasis([]float64{1.5, 0, -2.1}),
			scale:  10,
			wanted: NewBasis([]float64{15.0, 0, -21}),
		},
	}

	for i, tt := range tests {
		if got := tt.b.Scale(tt.scale); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.Scale(%f) = %v, want %v", i, tt.b, tt.scale, got, tt.wanted)
		}
	}
}

func TestBasisNorm(t *testing.T) {
	var tests = []struct {
		b      *Basis
		wanted float64
	}{
		{
			b:      NewBasis([]float64{1, 2, 2}),
			wanted: 3.0,
		},
	}

	for i, test := range tests {
		if got := test.b.Norm(); got != test.wanted {
			t.Errorf("test %d: %v.Norm() = %f, want %f", i, test.b, got, test.wanted)
		}
	}
}

func TestBasisSort(t *testing.T) {
	var tests = []struct {
		bs     BasisSlice
		wanted BasisSlice
	}{
		{
			bs: BasisSlice{
				NewBasis([]float64{1, 2, 2}),
				NewBasis([]float64{0, 0, 1}),
				NewBasis([]float64{0, -1, 1}),
			},
			wanted: BasisSlice{
				NewBasis([]float64{0, 0, 1}),
				NewBasis([]float64{0, -1, 1}),
				NewBasis([]float64{1, 2, 2}),
			},
		},
	}
	for i, tt := range tests {
		init := tt.bs
		sort.Sort(tt.bs)
		if !reflect.DeepEqual(tt.bs, tt.wanted) {
			t.Errorf("%d: sort(%v) got %v, want %v", i, init, tt.bs, tt.wanted)
		}
	}
}

func TestLatticeNew(t *testing.T) {
	var tests = []struct {
		l      []float64
		wanted *Lattice
	}{
		{
			l: []float64{1.1, 2.10, 3,
				1, 2, 3,
				1, 2, 3.000},
			wanted: &Lattice{*mat.NewDense(3, 3, []float64{1.1, 2.10, 3, 1, 2, 3, 1, 2, 3})},
		},
	}

	for i, test := range tests {
		got := NewLattice(test.l)
		if !reflect.DeepEqual(got, test.wanted) {
			t.Errorf("unexpected Lattice for test %d, got %v, want %v", i, got, test.wanted)
		}
	}
}

func TestLatticeDet(t *testing.T) {
	var tests = []struct {
		lat    *Lattice
		wanted float64
	}{
		{
			lat: NewLattice([]float64{6, 1, 1,
				4, -2, 5,
				2.0, 8, 7}),
			wanted: -306,
		},
		{
			lat: NewLattice([]float64{1, 0, 0,
				-0, -2, -0,
				1, 2, 0}),
			wanted: 0,
		},
	}

	var prec = 1e-5
	for i, test := range tests {
		if got := test.lat.Det(); !(math.Abs(got-test.wanted) < prec) {
			t.Errorf("%d: %v.Det() = %f, want %f", i, test.lat, got, test.wanted)
		}
	}
}

func TestLatticeT(t *testing.T) {
	var tests = []struct {
		lat    *Lattice
		wanted *Lattice
	}{
		{
			lat: NewLattice([]float64{6, 1, 1,
				4, -2, 5,
				2, 8, 7}),
			wanted: NewLattice([]float64{6, 4, 2,
				1, -2, 8,
				1, 5, 7}),
		},
	}
	for i, tt := range tests {
		if got := tt.lat.T(); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.T() = %v, want %v", i, tt.lat, got, tt.wanted)
		}
	}
}

func TestLatticeRotate(t *testing.T) {
	var tests = []struct {
		lat    *Lattice
		m      mat.Dense
		wanted *Lattice
	}{
		{
			NewLattice([]float64{6, 0, 0,
				0, 6, 0,
				0, 0, 6}),
			*mat.NewDense(3, 3, []float64{0, -1, 0, 0, 0, 1, -1, 0, 0}),
			NewLattice([]float64{0, -6, 0,
				0, 0, 6,
				-6, 0, 0}),
		},
	}

	for i, tt := range tests {
		got := tt.lat
		got.Rotate(tt.m)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.TransformBy(%v) got %v, want %v", i, tt.lat, tt.m.RawMatrix(), got, tt.wanted)
		}
	}
}

func TestLatticeGram(t *testing.T) {
	var tests = []struct {
		lat    *Lattice
		wanted *mat.Dense
	}{
		{
			lat: NewLattice([]float64{1, 2, 3,
				4, 5, 6,
				7, 8, 9}),
			wanted: mat.NewDense(3, 3, []float64{14, 32, 50, 32, 77, 122, 50, 122, 194}),
		},
	}

	for i, tt := range tests {
		got := tt.lat.Gram()
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.Gram() got %v(%T), want %v(%T)", i, tt.lat, got, got, tt.wanted, tt.wanted)
		}
	}
}

func TestLatticeSymmetry(t *testing.T) {
	var tests = []struct {
		latt   *Lattice
		wanted int
	}{
		{
			latt: NewLattice([]float64{0, 0.5, 0.5,
				0.5, 0, 0.5,
				0.5, 0.5, 0}),
			wanted: 48,
		},
		{
			latt: NewLattice([]float64{2.519, 0, 0,
				-1.2595, 2.181518, 0,
				0, 0, 4.091}),
			wanted: 24,
		},
		{
			latt: NewLattice([]float64{4, 0, 0,
				0, 4, 0,
				0, 0, 4}),
			wanted: 48,
		},
		{
			latt: NewLattice([]float64{4, 0, 0,
				4, 4, 0,
				0, 0, 11}),
			wanted: 8, // OR 8? symmetry of Lattice or symmetry of Lattice cell?
		},
	}

	for i, tt := range tests {
		ps := tt.latt.Symmetry(1e-5)
		if got := len(ps); got != tt.wanted {
			t.Errorf("%d: lattice %v, symmetry len is %d, want %d", i, tt.latt, got, tt.wanted)
		}
	}
}

func BenchmarkLatticeSymmetry(b *testing.B) {
	latt := NewLattice([]float64{2.519, 0, 0,
		-1.2595, 2.181518, 0,
		0, 0, 4.091})
	for i := 0; i < b.N; i++ {
		latt.Symmetry(1e-5)
	}
}

func TestNewFPos(t *testing.T) {
	var tests = []struct {
		data   []float64
		wanted *FPos
	}{
		{nil, &FPos{*mat.NewVecDense(3, nil)}},
		{[]float64{4, 5, 6}, &FPos{*mat.NewVecDense(3, []float64{4, 5, 6})}},
	}

	for i, tt := range tests {
		if got := NewFPos(tt.data); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: NewFPos(%v) got %v, want %v", i, tt.data, got, tt.wanted)
		}
	}
}

func TestFPosSub(t *testing.T) {
	var tests = []struct {
		fp     *FPos
		fpx    *FPos
		wanted *FPos
	}{
		{
			NewFPos([]float64{3, 2, 1}),
			NewFPos([]float64{1, 1, 0.2}),
			NewFPos([]float64{2, 1, 0.8}),
		},
	}

	for i, tt := range tests {
		if got := tt.fp.Sub(&tt.fpx.v); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.Sub(%v) got %v want %v", i, tt.fp, tt.fpx.v, got, tt.wanted)
		}
	}
}

func TestFPosAdd(t *testing.T) {
	var tests = []struct {
		fp     *FPos
		fpx    *FPos
		wanted *FPos
	}{
		{
			NewFPos([]float64{3, 2, 1}),
			NewFPos([]float64{1, 1, 0.2}),
			NewFPos([]float64{4, 3, 1.2}),
		},
	}

	for i, tt := range tests {
		if got := tt.fp.Add(&tt.fpx.v); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.Add(%v) got %v want %v", i, tt.fp, tt.fpx.v, got, tt.wanted)
		}
	}
}

func TestFPosEqual(t *testing.T) {
	var tests = []struct {
		fp     *FPos
		v      *FPos
		wanted bool
	}{
		{
			NewFPos([]float64{0.5, 0.25, 0.25}),
			NewFPos([]float64{0.5, 0.25, 0.25}),
			true,
		},
		{
			NewFPos([]float64{0.5, 0.25, 0.25}),
			NewFPos([]float64{0.5, 0.25, 0.75}),
			false,
		},
	}

	prec := 1e-5
	for i, tt := range tests {
		if got := tt.fp.Equal(tt.v, prec); got != tt.wanted {
			t.Errorf("%d: %v.Eqaul(%v) got %t", i, tt.fp, tt.v, got)
		}
	}
}

func TestFPosInCell(t *testing.T) {
	var tests = []struct {
		fp     *FPos
		wanted *FPos
	}{
		{
			NewFPos([]float64{0.5, -1.75, 0.25}),
			NewFPos([]float64{0.5, 0.25, 0.25}),
		},
		{
			NewFPos([]float64{0, 0, -1.00000001}),
			NewFPos([]float64{0, 0, 0}),
		},
	}

	prec := 1e-5
	for i, tt := range tests {
		if got := tt.fp.InCell(prec); !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: %v.InCell(%f) got %v, want %v", i, tt.fp, prec, got, tt.wanted)
		}
	}
}

func TestFPosEqualInCell(t *testing.T) {
	var tests = []struct {
		fp     *FPos
		fpx    *FPos
		wanted bool
	}{
		{
			NewFPos([]float64{0.5, -1.75, 0.25}),
			NewFPos([]float64{0.5, 0.25, 0.25}),
			true,
		},
	}

	prec := 1e-5
	for i, tt := range tests {
		if got := tt.fp.EqualInCell(tt.fpx, prec); got != tt.wanted {
			t.Errorf("%d: %v.EqaulInCell(%v) got %t", i, tt.fp, tt.fpx, got)
		}
	}
}

// func TestLeastEle(t *testing.T) {
// 	var tests = []struct {
// 		es     EType
// 		wanted string
// 	}{
// 		{[]string{"O", "O", "Li", "Li", "O", "O"}, "Li"},
// 		{[]string{"O", "O", "H"}, "H"},
// 	}
//
// 	for i, tt := range tests {
// 		if got := tt.es.leastEle(); got != tt.wanted {
// 			t.Errorf("%d: %v.leastEle() got %s want %s", i, tt.es, got, tt.wanted)
// 		}
// 	}
// }
//
// func TestIdxOfEle(t *testing.T) {
// 	var tests = []struct {
// 		es     EType
// 		s      string
// 		wanted []int
// 	}{
// 		{[]string{"O", "O", "Li", "Li", "O", "O"}, "Li", []int{2, 3}},
// 		{[]string{"O", "O", "H"}, "H", []int{2}},
// 		{[]string{"O", "O", "H"}, "O", []int{0, 1}},
// 	}
//
// 	for i, tt := range tests {
// 		if got := tt.es.IdxOfEle(tt.s); !reflect.DeepEqual(got, tt.wanted) {
// 			t.Errorf("%d: %v.IdxOfEle(%s) got %v want %v", i, tt.es, tt.s, got, tt.wanted)
// 		}
// 	}
// }

func TestCellSymmetry(t *testing.T) {
	var tests = []struct {
		c      *Cell
		wanted int
	}{
		{
			c: NewCell([]float64{3.8356, 0, 0,
				-1.918, 3.32207, 0,
				0, 0, 6.2770},
				[]float64{
					0.333333333, 0.666666667, 0,
					0.666666667, 0.333333333, 0.5,
					0.333333333, 0.666666667, 0.375,
					0.666666667, 0.333333333, 0.875},
				[]string{"Zn", "Zn", "S", "S"},
			),
			wanted: 12,
		},
		{
			c: NewCell([]float64{0, 2.8, 2.8,
				2.8, 0, 2.8,
				2.8, 2.8, 0},
				[]float64{
					0.75, 0.75, 0.75,
					0.0, 0.0, 0.0},
				[]string{"Zn", "S"},
			),
			wanted: 24,
		},
		{
			c: NewCell([]float64{5.4, 0, 0,
				0, 5.4, 0,
				0, 0, 5.4},
				[]float64{
					0.0, 0.0, 0.0,
					0.0, 0.5, 0.5,
					0.5, 0.0, 0.5,
					0.5, 0.5, 0.0,
					0.25, 0.25, 0.25,
					0.75, 0.75, 0.25,
					0.75, 0.25, 0.75,
					0.25, 0.75, 0.75},
				[]string{"Zn", "Zn", "Zn", "Zn", "S", "S", "S", "S"},
			),
			wanted: 96,
		},
		{
			c: NewCell([]float64{-2.877, 2.877, 5.1475,
				2.877, -2.877, 5.1475,
				2.877, 2.877, -5.1475},
				[]float64{
					0.0, 0.0, 0.0,
					0.75, 0.25, 0.5,
					0.25, 0.75, 0.5,
					0.5, 0.5, 0.0,
					0.179, 0.125, 0.554,
					0.375, 0.821, 0.946,
					0.571, 0.625, 0.446,
					0.875, 0.429, 0.054},
				[]string{"Ag", "Ag", "Ga", "Ga", "S", "S", "S", "S"},
			),
			wanted: 8,
		},
	}

	var prec = 1e-5
	for i, tt := range tests {
		got := len(tt.c.Symmetry(prec))
		if got != tt.wanted {
			t.Errorf("%d: Cell:\n %v symmetry length got %d, want %d", i, tt.c, got, tt.wanted)
			t.Errorf("%d: PointSym length %d", i, len(tt.c.Lattice.Symmetry(prec)))
			for _, v := range tt.c.Symmetry(prec) {
				t.Errorf("%d: pointOPS:%v\n transOPS%v\n", i, v.R, v.T)
			}
		}
	}
}

func TestCellFind(t *testing.T) {
	latt := []float64{5.4, 0, 0,
		0, 5.4, 0,
		0, 0, 5.4}
	frac := []float64{
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.5,
		0.5, 0.0, 0.5,
		0.5, 0.5, 0.0,
		0.25, 0.25, 0.25,
		0.75, 0.75, 0.25,
		0.75, 0.25, 0.75,
		0.25, 0.75, 0.75}
	et := []string{"Zn", "Zn", "Zn", "Zn", "S", "S", "S", "S"}
	ctest := NewCell(latt, frac, et)
	var tests = []struct {
		c        *Cell
		v        *FPos
		e        string
		wantedI  int
		wantedok bool
	}{
		{
			c:        ctest,
			v:        NewFPos([]float64{0.5, 0, 0.5}),
			e:        "Zn",
			wantedI:  2,
			wantedok: true,
		},
		{
			c:        ctest,
			v:        NewFPos([]float64{0.5, 0, 0.5}),
			e:        "S",
			wantedI:  0,
			wantedok: false,
		},
		{
			c:        ctest,
			v:        NewFPos([]float64{0.25, -2.25, 0.75}),
			e:        "S",
			wantedI:  7,
			wantedok: true,
		},
		{
			c:        ctest,
			v:        NewFPos([]float64{0., 0, -1.000000001}),
			e:        "Zn",
			wantedI:  0,
			wantedok: true,
		},
	}

	prec := 1e-5
	for i, tt := range tests {
		gotI, gotok := tt.c.Find(tt.v, tt.e, prec)
		if gotok != tt.wantedok {
			t.Errorf("%d: %v.Find (%v, %s) == %t", i, tt.c, tt.v, tt.e, gotok)
			break
		}
		if gotI != tt.wantedI {
			t.Errorf("%d: %v.Find (%v, %s) got Idx %d, want %d", i, tt.c, tt.v, tt.e, gotI, tt.wantedI)
		}
	}
}

func TestCellSimilar(t *testing.T) {
	lat := []float64{5.4, 0, 0,
		0, 5.4, 0,
		0, 0, 5.4}
	frac := []float64{
		0.0, 0.0, -1.000000001,
		0.0, 0.5, 0.25,
		0.5, 0.0, 0.5,
		0.5, 0.5, 0.0}
	fracx := []float64{
		0.0, 0.0, 0,
		0.5, -0.0, -0.5,
		0.0, -10.5, -0.75,
		0.5, 0.5, 0.0}

	var tests = []struct {
		c      *Cell
		cx     *Cell
		wanted bool
	}{
		{
			c:      NewCell(lat, frac, []string{"S", "S", "S", "S"}),
			cx:     NewCell(lat, fracx, []string{"S", "S", "S", "S"}),
			wanted: true,
		},
		{
			c:      NewCell(lat, frac, []string{"S", "S", "S", "S"}),
			cx:     NewCell(lat, fracx, []string{"S", "S", "Zn", "Zn"}),
			wanted: false,
		},
	}

	prec := 1e-5
	for i, tt := range tests {
		if got := tt.c.Equal(tt.cx, prec); got != tt.wanted {
			t.Errorf("%d: %v and %v are similar? [%t]", i, tt.c, tt.cx, got)
		}
	}
}
