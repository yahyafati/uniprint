package tools

// Dummy printer struct and list
type Printer struct {
	Name          string
	Colour        string
	PrintingSpeed string
	Model         string
}

var PRINTERS_LIST = []Printer{
	{"tfppr1", "b/w", "35 p/m", "Ricoh MP 3554"},
	{"tfppr2", "Color", "38 p/m", "HP Color LaserJet Enterprise M553"},
	{"tfppr3", "b/w", "45 p/m", "HP LaserJet Enterprise MFP M528"},
	{"tfppr4", "Color", "56 p/m", "HP Color LaserJet Enterprise M653"},
}
