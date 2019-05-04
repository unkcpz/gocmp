package crystal

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNewPointOp(t *testing.T) {
	var tests = []struct {
		data   [][]float64
		wanted *PointOp
	}{
		{
			data: [][]float64{[]float64{1, 0, 1}, []float64{-1, 0, 1}, []float64{0, 1, -1}},
			wanted: &PointOp{
				*mat.NewDense(3, 3, []float64{1, 0, 1, -1, 0, 1, 0, 1, -1}),
			},
		},
	}

	for i, tt := range tests {
		got := NewPointOp(tt.data)
		if !reflect.DeepEqual(got, tt.wanted) {
			t.Errorf("%d: unexpected NewPoint(%v), got %v, want %v", i, tt.data, got, tt.wanted)
		}
	}
}

// func TestPointOpActOn(t *testing.T) {
// 	var tests = []struct {
// 		r      *PointOp
// 		input  *Fractions
// 		wanted *Fractions
// 	}{
// 		{
// 			NewPointOp([][]float64{
// 				[]float64{0, 1, 0},
// 				[]float64{-1, -1, 0},
// 				[]float64{0, 0, 1},
// 			}),
// 			NewFractions(2, []float64{0, 0, 0,
// 				0.25, 0.25, 0.25}),
// 			NewFractions(2, []float64{0, 0, 0, -0.25,
// 				0, 0.25}),
// 		},
// 	}
//
// 	for i, tt := range tests {
// 		got := tt.input
// 		tt.r.ActOn(got)
// 		if !reflect.DeepEqual(got, tt.wanted) {
// 			t.Errorf("%d: %v.Operate(%v) got %v, want %v", i, tt.input, tt.r, got, tt.wanted)
// 		}
// 	}
// }

func TestDupHnfs(t *testing.T) {
	var tests = []struct {
		input  int
		wanted int
	}{
		{2, 7},
		{6, 93},
	}

	for i, tt := range tests {
		if got := len(DupHnfs(tt.input)); got != tt.wanted {
			t.Errorf("%d: hnfs(%d) length %d, want %d", i, tt.input, got, tt.wanted)
		}
	}
}

func BenchmarkDupHnfs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DupHnfs(6)
	}
}
