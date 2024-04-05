package tflexregex

import "errors"

type ast struct {
	elems    []elem
	parentOf map[uint][]uint
	childOf  map[uint]uint
}

const (
	codeNone byte = iota
	codeEscapeLiteral
	codeDot
	codeDecimal
	codeLiteral
	codeClass
	codeEscape
	codeTerm
	codeModifier
	codeParPattern
	codePattern
	codeModified
	codeExpression
	codeRegex
)

type elem struct {
	code    byte
	content []byte
}

func newAst() ast {
	return ast{
		elems:    make([]elem, 1, 256),
		parentOf: make(map[uint][]uint, 256),
		childOf:  make(map[uint]uint, 256),
	}
}

func makeRoot(tree ast, code byte, content []byte) ast {
	tree.elems[0] = elem{
		code:    code,
		content: content,
	}

	return tree
}

func addChild(tree ast, parent uint, code byte, content []byte) (ast, error) {
	length := uint(len(tree.elems))
	if parent >= length {
		return tree, errors.New("invalid parent index")
	}

	tree.elems = append(tree.elems, elem{
		code:    code,
		content: content,
	})

	tree.parentOf[parent] = append(tree.parentOf[parent], length)
	tree.parentOf[length] = make([]uint, 0, 8)
	tree.childOf[length] = parent
	return tree, nil
}
