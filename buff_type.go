package main

type BuffType int

const (
	NONE BuffType = iota
	BLUE
	SB
)

var ClassToBuffMapping = map[string]BuffType{
	"SIN": BLUE,
	"KFM": BLUE,
	"WAR": SB,
	"WL":  SB,
}

var BuffToStringMapping = map[BuffType]string{
	BLUE: "BB",
	SB:   "SB",
	NONE: "  ",
}
