package crystal

// func TestDealaunayReduce(t *testing.T) {
// 	var prec = 1e-5
// 	var tests = []struct {
// 		oLatt  *Lattice
// 		wanted *Lattice
// 	}{
// 		{
// 			oLatt: NewLattice([]float64{1, 0.0, 0.0,
// 				0, 2, 0.0,
// 				0, 0, 3}),
// 			wanted: NewLattice([]float64{-1, 0, 0,
// 				0, 2, 0,
// 				0, 0, -3}),
// 		},
// 		{
// 			oLatt: NewLattice([]float64{4.693, 0.0, 0.0,
// 				-0.0577, 4.93566, 0.0,
// 				0.05647, -4.93566, 5.678}),
// 			wanted: NewLattice([]float64{-4.693, 0, 0,
// 				-0.0577, 4.93566, 0,
// 				0.001230000000000002, 0, -5.678}),
// 		},
// 		// {
// 		// 	oLatt: &Lattice{Basis{[]float64{2.519, 0.0, 0.0}},
// 		// 		Basis{[]float64{-1.2595, 2.181518, 0.0}},
// 		// 		Basis{[]float64{0.0, 0.0, 4.091}}},
// 		// 	wanted: &Lattice{Basis{[]float64{2.519, 0.0, 0.0}},
// 		// 		Basis{[]float64{-1.2595, 2.181518, 0.0}},
// 		// 		Basis{[]float64{0.0, 0.0, 4.091}}},
// 		// }, //bugs should be a tri-octuse but a tri-acute
// 	}
//
// 	for i, test := range tests {
// 		got, err := DealaunayReduce(test.oLatt, prec)
// 		if err != nil {
// 			t.Errorf("%d: can't DealaunayReduce(%v, %f)", i, test.oLatt, prec)
// 		}
// 		if !reflect.DeepEqual(got, test.wanted) {
// 			t.Errorf("%d: DealaunayReduce(%v, %f) = %v, want %v", i, test.oLatt, prec, got, test.wanted)
// 		}
// 	}
// }
