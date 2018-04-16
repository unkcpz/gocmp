package vaspxml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const xmldata = `<?xml version="1.0" encoding="ISO-8859-1"?>
<modeling>
 <generator>
  <i name="program" type="string">vasp </i>
  <i name="version" type="string">5.4.4.18Apr17-6-g9f103f2a35  </i>
  <i name="subversion" type="string">(build Nov 17 2017 16:18:39) complex            parallel </i>
  <i name="platform" type="string">LinuxIFC </i>
  <i name="date" type="string">2017 12 27 </i>
  <i name="time" type="string">14:18:28 </i>
 </generator>
 <atominfo>
   <atoms>       5 </atoms>
   <types>       3 </types>
   <array name="atoms" >
    <dimension dim="1">ion</dimension>
    <field type="string">element</field>
    <field type="int">atomtype</field>
    <set>
     <rc><c>Sc</c><c>   1</c></rc>
     <rc><c>Sc</c><c>   1</c></rc>
     <rc><c>Sc</c><c>   1</c></rc>
     <rc><c>P </c><c>   2</c></rc>
     <rc><c>O </c><c>   3</c></rc>
    </set>
   </array>
  </atominfo>
 <structure name="initialpos" >
  <crystal>
   <varray name="basis" >
    <v>       2.71500000       2.71500000       0.00000000 </v>
    <v>       0.00000000       2.71500000       2.71500000 </v>
    <v>       2.71500000       0.00000000       2.71500000 </v>
   </varray>
   <i name="volume">     40.02575175 </i>
   <varray name="rec_basis" >
    <v>       0.18416206       0.18416206      -0.18416206 </v>
    <v>      -0.18416206       0.18416206       0.18416206 </v>
    <v>       0.18416206      -0.18416206       0.18416206 </v>
   </varray>
  </crystal>
  <varray name="positions" >
   <v>       0.00000000       0.00000000       0.00000000 </v>
   <v>       0.25000000       0.25000000       0.25000000 </v>
  </varray>
 </structure>
 <structure>
  <crystal>
   <varray name="basis" >
    <v>       2.71500000       2.71500000       0.00000000 </v>
    <v>       0.00000000       2.71500000       2.71500000 </v>
    <v>       2.71500000       0.00000000       2.71500000 </v>
   </varray>
   <i name="volume">     40.02575175 </i>
   <varray name="rec_basis" >
    <v>       0.18416206       0.18416206      -0.18416206 </v>
    <v>      -0.18416206       0.18416206       0.18416206 </v>
    <v>       0.18416206      -0.18416206       0.18416206 </v>
   </varray>
  </crystal>
  <varray name="positions" >
   <v>       0.00000000       0.00000000       0.00000000 </v>
   <v>       0.25000000       0.25000000       0.25000000 </v>
  </varray>
 </structure>
 <structure name="finalpos" >
  <crystal>
   <varray name="basis" >
    <v>       2.71500000       2.71500000       0.00000000 </v>
    <v>       0.00000000       2.71500000       2.71500000 </v>
    <v>       2.71500000       0.00000000       2.71500000 </v>
   </varray>
   <i name="volume">     40.02575175 </i>
   <varray name="rec_basis" >
    <v>       0.18416206       0.18416206      -0.18416206 </v>
    <v>      -0.18416206       0.18416206       0.18416206 </v>
    <v>       0.18416206      -0.18416206       0.18416206 </v>
   </varray>
  </crystal>
  <varray name="positions" >
   <v>       0.00000000       0.00000000       0.00000000 </v>
   <v>       0.25000000       0.25000000       0.25000000 </v>
  </varray>
 </structure>
</modeling>
`

func twoDSliceEqual(a, b [][]float64) bool {
	if a == nil && b == nil {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for i, vo := range a {
		if len(vo) != len(b[i]) {
			return false
		}
		for j, vi := range vo {
			if b[i][j] != vi {
				return false
			}
		}
	}
	return true
}

func TestParseCalInfo(t *testing.T) {
	vasprun, _ := Parse(strings.NewReader(xmldata))
	tags := make(map[string]string)
	for _, tag := range vasprun.CalInfo.Tags {
		tags[tag.Key] = tag.Value
	}

	var tests = []struct {
		input  string
		wanted string
	}{
		{"program", "vasp "},
		{"platform", "LinuxIFC "},
		{"date", "2017 12 27 "},
	}

	for _, test := range tests {
		if got := tags[test.input]; got != test.wanted {
			t.Errorf("Generator Key: %s, Value %s, want %s", test.input, got, test.wanted)
		}
	}
}

func TestParseAtomInfo(t *testing.T) {
	vasprun, _ := Parse(strings.NewReader(xmldata))
	// test []element
	var tests = []struct {
		input  int
		wanted string
	}{
		{0, "Sc"},
		{1, "Sc"},
		{2, "Sc"},
		{3, "P "},
		{4, "O "},
	}
	for _, test := range tests {
		if got := vasprun.AtomInfo.Elements[test.input]; got != test.wanted {
			t.Errorf("vasprun.AtomInfo.elements[%d] == %v, want %v",
				test.input, got, test.wanted)
		}
	}
}

func TestParseFinalStructure(t *testing.T) {
	vasprun, err := Parse(strings.NewReader(xmldata))
	if err != nil {
		t.Error(err)
	}

	wanted := [][]float64{
		{2.715, 2.715, 0.000},
		{0.000, 2.715, 2.715},
		{2.715, 0.000, 2.715},
	}
	if got := vasprun.FinalStructure.Lattice; !twoDSliceEqual(got, wanted) {
		t.Errorf("vasprun.FinalStructure.Lattice == %v, want %v", got, wanted)
	}

	wanted = [][]float64{
		{0.00, 0.00, 0.00},
		{0.25, 0.25, 0.25},
	}
	if got := vasprun.FinalStructure.Positions; !twoDSliceEqual(got, wanted) {
		t.Errorf("vasprun.FinalStructure.Positions == %v, want %v", got, wanted)
	}
}

func TestParseInitialStructure(t *testing.T) {
	vasprun, err := Parse(strings.NewReader(xmldata))
	if err != nil {
		t.Error(err)
	}

	wanted := [][]float64{
		{2.715, 2.715, 0.000},
		{0.000, 2.715, 2.715},
		{2.715, 0.000, 2.715},
	}
	if got := vasprun.InitialStructure.Lattice; !twoDSliceEqual(got, wanted) {
		t.Errorf("vasprun.InitialStructure.Lattice == %v, want %v", got, wanted)
	}

	wanted = [][]float64{
		{0.00, 0.00, 0.00},
		{0.25, 0.25, 0.25},
	}
	if got := vasprun.InitialStructure.Positions; !twoDSliceEqual(got, wanted) {
		t.Errorf("vasprun.InitialStructure.Positions == %v, want %v", got, wanted)
	}
}

func TestParseStructures(t *testing.T) {
	vasprun, _ := Parse(strings.NewReader(xmldata))
	wanted := 3
	if got := len(vasprun.Structures); got != wanted {
		t.Errorf("len(vasprun.Structures) == %d, want %d", got, wanted)
	}
}

func TestParseFromFile(t *testing.T) {
	absFilePath, err := filepath.Abs("test.xml")
	if err != nil {
		t.Error(err)
	}
	xmlfile, err := os.Open(absFilePath)
	if err != nil {
		t.Error(err)
	}
	defer xmlfile.Close()

	vasprun, err := Parse(xmlfile)
	if err != nil {
		t.Error(err)
	}
	wanted := 4
	if got := len(vasprun.Structures); got != wanted {
		t.Errorf("parsing Si-bs structures xml got %d structures, want %d", got, wanted)
	}
}

func TestParseDOS(t *testing.T) {
	absFilePath, err := filepath.Abs("test.xml")
	if err != nil {
		t.Error(err)
	}
	xmlfile, err := os.Open(absFilePath)
	if err != nil {
		t.Error(err)
	}
	defer xmlfile.Close()

	vasprun, err := Parse(xmlfile)
	if err != nil {
		t.Error(err)
	}

	wanted := 301
	dos := vasprun.DOS
	if got := len(dos.X); got != wanted {
		t.Errorf(`DOS.X len(%v) == %d, want %d`, dos.X, got, wanted)
	}
	if got := len(dos.TDOS["up"]); got != wanted {
		t.Errorf(`DOS.TDOS["up"] len(%v) == %d, want %d`, dos.TDOS, got, wanted)
	}
	if got := len(dos.IntDOS["down"]); got != wanted {
		t.Errorf(`DOS.IntDOS["down"] len(%v) == %d, want %d`, dos.IntDOS, got, wanted)
	}

	pdos := vasprun.DOS.PDOS
	wanted = 36 // 2(atoms)*9(orbital)*2(spins) = 36 states
	if got := len(pdos); got != wanted {
		t.Errorf(`DOS.PDOS length == %d, want %d`, got, wanted)
	}
}

func TestTry(t *testing.T) {

}
