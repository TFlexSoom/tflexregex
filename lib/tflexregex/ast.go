package tflexregex

type ast struct {
	root      uint
	elem      []elem
	relations map[uint]uint
}

const (
	codeNone byte = iota
	codeLiteral
	codeClass
	codeDot
	codeUnion
	codeSub
)

type elem struct {
	code    byte
	lte     uint
	gte     uint
	content []rune
}

func newLiteralElem(v rune) elem {
	return elem{
		code:    codeLiteral,
		lte:     1,
		gte:     1,
		content: []rune{v},
	}
}

func newClassElem(runes []rune) elem {
	return elem{
		code:    codeClass,
		lte:     1,
		gte:     1,
		content: runes,
	}
}

func newDotElem() elem {
	return elem{
		code:    codeDot,
		lte:     1,
		gte:     1,
		content: []rune{},
	}
}

func newUnionElem() elem {
	return elem{
		code:    codeUnion,
		lte:     1,
		gte:     1,
		content: []rune{},
	}
}

func newSubElem(subRef uint) elem {
	return elem{
		code:    codeSub,
		lte:     1,
		gte:     1,
		content: []rune{},
	}
}

func newAst() ast {
	return ast{
		root:      0,
		elem:      make([]elem, 0, 256),
		relations: make(map[uint]uint, 256),
	}
}
