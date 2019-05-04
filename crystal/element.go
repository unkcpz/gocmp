package crystal

func SymToNum(s string) int {
  return Element[s]
}

func NumToSym(a int) string {
  return Symbol[a]
}

var Element = map[string]int {
  "H" : 1,
}

var Symbol = map[int]string {
  1 : "H",
}
