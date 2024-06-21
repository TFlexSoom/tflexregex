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