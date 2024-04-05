package tflexregex

import (
	"errors"
	"fmt"
)

/*
REGEX GRAMMAR

REGEX => EXPRESSION | EXPRESSION REGEX
EXPRESSION => MODIFIED | MODIFIED '|' MODIFIED
MODIFIED => PATTERN MODIFIER?
PATTERN => '(' PAR_PATTERN ')' | TERM
PAR_PATTERN => '(' PAR_PATTERN ')' | PAR_PATTERN '|' PAR_PATTERN | TERM+
MODIFIER => {DECIMAL,DECIMAL} | {DECIMAL,} | {,DECIMAL} | '+' | '*'
TERM => CLASS | LITERAL | ESCAPE | DECIMAL | DOT
ESCAPE => ESCAPE_LITERAL .
CLASS => '[' [^]]* ']'
LITERAL => [^123456789()'\'+*{}[\]]
DECIMAL => 123456789+
DOT => '.'
ESCAPE_LITERAL => '\'
*/

type parsingMonad struct {
	result ast
	index  uint
	input  string
	err    error
}

func newParsingMonad(input string) parsingMonad {
	return parsingMonad{
		result: newAst(),
		index:  0,
		input:  input,
		err:    nil,
	}
}

func parse(regex string) (ast, error) {
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
	pMonad.result
	length := uint(len(pMonad.input))
	for pMonad.index < length {
		pMonad = mapErr(pMonad, expression)
		if pMonad.err != nil {
			return pMonad
		}
	}

	return pMonad
}

func expression(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		pMonad.err = errors.New("expected expression, but end of input was reached")
		return pMonad
	}

	pMonad = mapErr(pMonad, pattern)

	if hasChar(pMonad, '|') {
		pMonad = mapErr(pMonad, bind(consumeChar, '|'))
		pMonad = mapErr(pMonad, expression)
	}

	return pMonad
}

func pattern(pMonad parsingMonad) parsingMonad {
	length := uint(len(pMonad.input))
	if pMonad.index >= length {
		pMonad.err = errors.New("expected expression, but end of input was reached")
		return pMonad
	}

	if hasChar(pMonad, '(') {
		pMonad = mapErr(pMonad, bind(consumeChar, '('))
		pMonad = mapErr(pMonad, pattern)
		pMonad = mapErr(pMonad, bind(consumeChar, ')'))
		return pMonad
	}

}
