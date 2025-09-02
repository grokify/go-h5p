package schemas

import (
	_ "embed"
)

//go:embed essay_semantics.json
var EssaySemanticsBytes []byte

//go:embed multichoice_semantics.json
var MultiChoiceSemanticsBytes []byte

//go:embed truefalse_semantics.json
var TrueFalseSemanticsBytes []byte
