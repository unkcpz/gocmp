// vaspxml parses vasprun.xml
package vaspxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"mellium.im/xmlstream"
)

type VaspRunXML struct {
	XMLName          xml.Name `xml:"modeling"`
	CalInfo          CalInfo
	AtomInfo         AtomInfo
	FinalStructure   Structure
	InitialStructure Structure
	Structures       []Structure
	DOS              DOS
}

type CalInfo struct {
	XMLName xml.Name `xml:"generator"`
	Tags    []Tag    `xml:"i"`
}

type Tag struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type AtomInfo struct {
	Elements []string
}

type Structure struct {
	Lattice   [][]float64
	Positions [][]float64
}

type DOS struct {
	X      []float64
	TDOS   map[Spin][]float64
	IntDOS map[Spin][]float64
	PDOS   map[State][]float64
}

type Spin string

type State struct {
	Ion     int
	Orbital string
	Spin    Spin
}

func Parse(r io.Reader) (VaspRunXML, error) {
	var vr VaspRunXML
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return vr, err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "generator" {
				var c CalInfo
				decoder.DecodeElement(&c, &se)
				vr.CalInfo = c
			}
			if se.Name.Local == "atominfo" {
				dosDec := Inner(decoder)
				vr.AtomInfo = handlerAtomInfo(dosDec)
			}
			if se.Name.Local == "structure" {
				dosDec := Inner(decoder)
				st, err := handlerStructrue(dosDec)
				if err != nil {
					return vr, err
				}
				vr.Structures = append(vr.Structures, st)

				vr.InitialStructure = vr.Structures[0]
				vr.FinalStructure = vr.Structures[len(vr.Structures)-1]
			}
			if se.Name.Local == "dos" {
				dosDec := xmlstream.Inner(decoder)
				vr.DOS = handlerDOS(dosDec)
			}
		}
	}
	return vr, nil
}

func handlerAtomInfo(tr xml.TokenReader) (ai AtomInfo) {
	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		decoder := xml.NewTokenDecoder(tr)
		if se, ok := tok.(xml.StartElement); ok {
			if len(se.Attr) > 0 && se.Attr[0].Value == "atoms" {
				type Atom struct {
					Name []string `xml:"c"`
				}
				type atominfo struct {
					Atoms []Atom `xml:"set>rc"`
				}
				var atoms atominfo
				decoder.DecodeElement(&atoms, &se)
				s := make([]string, len(atoms.Atoms))
				for i, a := range atoms.Atoms {
					s[i] = strings.Fields(a.Name[0])[0]
				}
				ai.Elements = s
			}
		}
	}
	return
}

func handlerDOS(tr xml.TokenReader) (dos DOS) {
	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		if se, ok := tok.(xml.StartElement); ok {
			if se.Name.Local == "total" {
				dosDec := xmlstream.Inner(tr)
				dos.X, dos.TDOS, dos.IntDOS = handlerTDOS(dosDec)
			}
			if se.Name.Local == "partial" {
				dosDec := xmlstream.Inner(tr)
				dos.PDOS = handlerPDOS(dosDec)
			}
		}
	}
	return dos
}

func handlerStructrue(tr xml.TokenReader) (st Structure, err error) {
	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			err = fmt.Errorf("handlerStructrue: dec.Token(): %v", err)
			return st, err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "varray" && se.Attr[0].Value == "basis" {
				st.Lattice, err = vParser(tr, &se)
				if err != nil {
					return st, fmt.Errorf("parsing xml lattice %v: %v", se, err)
				}
			}
			if se.Name.Local == "varray" && se.Attr[0].Value == "positions" {
				st.Positions, err = vParser(tr, &se)
				if err != nil {
					return st, fmt.Errorf("parsing xml position %v: %v", se, err)
				}
			}
		}
	}
	return st, nil
}

func handlerTDOS(tr xml.TokenReader) (x []float64, tdos, idos map[Spin][]float64) {
	tdos = make(map[Spin][]float64)
	idos = make(map[Spin][]float64)
	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		switch te := tok.(type) {
		case xml.StartElement:
			if te.Name.Local == "set" {
				if len(te.Attr) > 0 && te.Attr[0].Value == "spin 1" {
					var s Spin = "up"
					all, _ := rParser(tr, &te)
					x = make([]float64, len(all))
					tdos[s] = make([]float64, len(all))
					idos[s] = make([]float64, len(all))
					for i, row := range all {
						x[i], tdos[s][i], idos[s][i] = row[0], row[1], row[2]
					}
				}
				if len(te.Attr) > 0 && te.Attr[0].Value == "spin 2" {
					var s Spin = "down"
					all, _ := rParser(tr, &te)
					tdos[s] = make([]float64, len(all))
					idos[s] = make([]float64, len(all))
					for i, row := range all {
						_, tdos[s][i], idos[s][i] = row[0], row[1], row[2]
					}
				}
			}
		}
	}
	return
}

func handlerPDOS(tr xml.TokenReader) (pdos map[State][]float64) {
	pdos = make(map[State][]float64)
	ion := 0
	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		switch te := tok.(type) {
		case xml.StartElement:
			if te.Name.Local == "set" {
				if len(te.Attr) > 0 && strings.HasPrefix(te.Attr[0].Value, "ion") {
					ion, _ = strconv.Atoi(strings.Split(te.Attr[0].Value, " ")[1])
					ion = ion - 1
				}
				if len(te.Attr) > 0 && te.Attr[0].Value == "spin 1" {
					var s Spin = "up"
					atompdos, _ := rParser(tr, &te)
					splitOrbital(pdos, ion, s, atompdos)
				}
				if len(te.Attr) > 0 && te.Attr[0].Value == "spin 2" {
					var s Spin = "down"
					atompdos, _ := rParser(tr, &te)
					splitOrbital(pdos, ion, s, atompdos)
				}
			}
		}
	}
	return
}

func splitOrbital(pdos map[State][]float64, ion int, s Spin, ispdos [][]float64) {
	orbital := []string{"s", "py", "pz", "px", "dxy", "dyz", "dz2", "dxz", "x2-y2"}
	for i, o := range orbital {
		var state = State{Ion: ion, Orbital: o, Spin: s}
		pdos[state] = ispdos[i]
	}
}

func rParser(tr xml.TokenReader, te *xml.StartElement) ([][]float64, error) {
	type rowString struct {
		Row []string `xml:"r"`
	}
	var rs rowString
	dec := xml.NewTokenDecoder(tr)
	dec.DecodeElement(&rs, te)
	return vecRowParsing(rs.Row)
}

// vecParser parses vasprun.xml vector object into a two-dimenssional slice
func vParser(tr xml.TokenReader, se *xml.StartElement) ([][]float64, error) {
	type rowString struct {
		Row []string `xml:"v"`
	}
	var rs rowString
	dec := xml.NewTokenDecoder(tr)
	dec.DecodeElement(&rs, se)
	return vecRowParsing(rs.Row)
}

func vecRowParsing(rs []string) (v [][]float64, err error) {
	v = make([][]float64, len(rs))
	for i, arow := range rs {
		valueStrings := strings.Fields(arow)
		v[i] = make([]float64, len(valueStrings))
		for j, vitem := range valueStrings {
			var err error
			v[i][j], err = strconv.ParseFloat(vitem, 64)
			if err != nil {
				return nil, fmt.Errorf("rowParse %v: parse %v as float64: %v",
					arow, vitem, err)
			}
		}
	}
	return v, nil
}
