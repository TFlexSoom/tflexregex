package tflexregex

/*
REGEX GRAMMAR

REGEX => EXPRESSION MODIFIER? EXPRESSION*
EXPRESSION => ( EXPRESSION ) | EXPRESSION '|' EXPRESSION | CLASS | TERM+
MODIFIER => {DECIMAL,DECIMAL} | {DECIMAL,} | {,DECIMAL} | '+' | '*'
TERM => LITERAL | ESCAPE | DECIMAL
ESCAPE => ESCAPE_LITERAL .
CLASS => '[' [^]]* ']'
LITERAL => [^123456789()'\'+*{}[\]]
DECIMAL => 123456789+
ESCAPE_LITERAL => '\'
*/

type parsingMonad struct {
	result ast
	index  uint
}

func newParsingMonad() parsingMonad {
	return parsingMonad{
		result: newAst(),
		index:  0,
	}
}

func parse(regex string) (ast, error) {
	pMonad := newParsingMonad()
	for i, v := range regex {
		pMonad, err = descent(pMonad, v)

		if err != nil {
			return pMonad.result, err
		}
	}

	return pMonad.result
}

func descent(pMonad parsingMonad, v rune) (ast, error) {
	pMonad = expression(pMonad, v)

}
