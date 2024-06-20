package tflexregex

import (
	"errors"
	"fmt"

	tree "github.com/tflexsoom/go-tree/lib"
)

/*
REGEX GRAMMAR

REGEX => EXPRESSION REGEX?
EXPRESSION => PATTERN ('|' PATTERN)*
PATTERN => '(' EXPRESSION ')' MODIFIER?

	| TERM MODIFIER?

MODIFIER => {DECIMAL,DECIMAL} | {DECIMAL,} | {,DECIMAL} | '+' | '*'
TERM => ESCAPE | CLASS | DOT | LITERAL
ESCAPE => ESCAPE_LITERAL .
CLASS => '[' [^]]* ']'
LITERAL => [^()'\'+*{}[\]]
DECIMAL => 123456789+
DOT => '.'
ESCAPE_LITERAL => '\'
*/
type parsingMonad struct {
	index uint
	input string
	err   error
}

func newParsingMonad(input string) parsingMonad {
	tree := tree.NewGraphTreeCap[astElem](16, 4)
	return parsingMonad{
		result:   tree,
		lastTree: tree,
		index:    0,
		input:    input,
		err:      nil,
	}
}

func parse(regex string) (tree.Tree[astElem], error) {
	pMonad := newParsingMonad(regex)
	pMonad = root(pMonad)
	return pMonad.result, pMonad.err
}

type parsingMonadMapper func(parsingMonad) parsingMonad

func mapErr(pMonad parsingMonad, mapper parsingMonadMapper) parsingMonad {
	if pMonad.err != nil {
		return pMonad
	}

	return mapper(pMonad)
}

func bind[V any](mapping func(parsingMonad, V) parsingMonad, val V) func(parsingMonad) parsingMonad {
	return func(pMonad parsingMonad) parsingMonad {
		return mapping(pMonad, val)
	}
}

func has(pMonad parsingMonad, v rune) bool {
	val := uint32(v)
	for shiftVal := 8 * 3; shiftVal >= 0; shiftVal -= 8 {
		if !hasChar(pMonad, byte(val>>shiftVal)) {
			return false
		}
	}

	return true
}

func consume(pMonad parsingMonad, v rune) parsingMonad {
	val := uint32(v)
	for shiftVal := 8 * 3; shiftVal >= 0; shiftVal -= 8 {
		pMonad = mapErr(pMonad, bind(consumeChar, byte(val>>shiftVal)))
	}

	return pMonad
}

func hasChar(pMonad parsingMonad, v byte) bool {
	return pMonad.input[pMonad.index] == v
}

func consumeChar(pMonad parsingMonad, v byte) parsingMonad {
	nextInput := pMonad.input[pMonad.index]
	if nextInput != v {
		pMonad.err = fmt.Errorf("expected byte %v but got %v instead", v, nextInput)
		return pMonad
	}

	pMonad.index++
	return pMonad
}

func root(pMonad parsingMonad) parsingMonad {
	pMonad.result.SetValue(astElem{
		code:    codeRegex,
		content: []byte{},
	})

	length := uint(len(pMonad.input))
	for pMonad.index < length && pMonad.err == nil {
		pMonad = mapErr(pMonad, expression)
	}

	return pMonad
}

func expression(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		pMonad.err = errors.New("expected expression, but end of input was reached")
		return pMonad
	}

	expressionRoot := pMonad.lastTree.AddChild(astElem{
		code:    codeExpression,
		content: []byte{},
	})
	pMonad.lastTree = expressionRoot
	pMonad = mapErr(pMonad, pattern)

	for hasChar(pMonad, '|') && pMonad.err == nil {
		pMonad = mapErr(pMonad, bind(consumeChar, '|'))

		pMonad.lastTree = expressionRoot
		pMonad = mapErr(pMonad, pattern)
	}

	return pMonad
}

func pattern(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		pMonad.err = errors.New("expected pattern, but end of input was reached")
		return pMonad
	}

	patternRoot := pMonad.lastTree.AddChild(astElem{
		code:    codePattern,
		content: []byte{},
	})

	pMonad.lastTree = patternRoot
	if hasChar(pMonad, '(') {
		pMonad = mapErr(pMonad, bind(consumeChar, '('))
		pMonad = mapErr(pMonad, expression)
		pMonad = mapErr(pMonad, bind(consumeChar, ')'))
	} else {
		pMonad = mapErr(pMonad, term)
	}

	pMonad.lastTree = patternRoot
	return mapErr(pMonad, modifier)
}

// optionally added
func modifier(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		return pMonad
	}

	if hasChar(pMonad, '{') {
		// TODO
	} else if hasChar(pMonad, '+') {
		pMonad = mapErr(pMonad, bind(consumeChar, '+'))
		pMonad.lastTree.AddChild(astElem{
			code:    codeModifier,
			content: []byte{'+'},
		})
	} else if hasChar(pMonad, '*') {
		pMonad = mapErr(pMonad, bind(consumeChar, '+'))
		pMonad.lastTree.AddChild(astElem{
			code:    codeModifier,
			content: []byte{'*'},
		})
	}

	return pMonad
}

func term(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		pMonad.err = errors.New("expected term, but end of input was reached")
		return pMonad
	}

	termRoot := pMonad.lastTree.AddChild(astElem{
		code:    codeTerm,
		content: []byte{},
	})

	pMonad.lastTree = termRoot

	if hasChar(pMonad, '\\') {
		pMonad = mapErr(pMonad, escape)
	} else if hasChar(pMonad, '[') {
		pMonad = mapErr(pMonad, class)
	} else if hasChar(pMonad, '.') {
		pMonad = mapErr(pMonad, bind(consumeChar, '.'))
		pMonad.lastTree.AddChild(astElem{
			code:    codeDot,
			content: []byte{},
		})
	} else {
		pMonad = mapErr(pMonad, literal)
	}

	return pMonad
}

func escape(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index > length { // should be >= length + 1 because \\ should prefix another character
		pMonad.err = errors.New("expected escape, but end of input was reached")
		return pMonad
	}

	pMonad = mapErr(pMonad, bind(consumeChar, '\\'))
	nextChar := pMonad.input[pMonad.index]
	pMonad.index += 1
	pMonad.lastTree.AddChild(astElem{
		code:    codeEscape,
		content: []byte{nextChar},
	})

	return pMonad
}

func class(pMonad parsingMonad) parsingMonad {

}
