package crystal

import (
  "fmt"

	"gonum.org/v1/gonum/mat"
)

type Rotation struct {
  data *mat.Dense
}

func (r Rotation) String() string {
  fa := mat.Formatted(r.data, mat.Squeeze())
  return fmt.Sprintf("%v", fa)
}

type Translation struct {
  data *mat.VecDense
}

func (t Translation) String() string {
  fa := mat.Formatted(t.data.T(), mat.Squeeze())
  return fmt.Sprintf("%v", fa)
}
