package parse

import "fmt"

type parsingOutput struct{}

func parseRegex(regex string) parsingOutput {
	return parsingOutput{}
}

func parse(regex string) (, error) {
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
