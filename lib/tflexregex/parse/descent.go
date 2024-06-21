package parse

/*
REGEX GRAMMAR

REGEX => '^'? EXPRESSION? '$'?
EXPRESSION => MODIFIED_TERM + ('|' EXPRESSION )*
MODIFIED_TERM => '(' EXPRESSION ')' MODIFIER?
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

type parsingMonad interface {
	isEmpty() bool
	has(byte) bool
	within(byte, byte) bool
	errorOut(string)
	acceptOnHas(byte)
	accept(byte)
	acceptByte()
	acceptWithin(byte, byte)
	acceptUnicode()
	ready(func(parsingMonad) parsingMonad)
}

func regex(pMonad parsingMonad) parsingMonad {
	pMonad.acceptOnHas('^')
	pMonad.ready(expression)
	pMonad.ready(func(pm parsingMonad) parsingMonad {
		pm.acceptOnHas('$')
		return pm
	})
	return pMonad
}

func expression(pMonad parsingMonad) parsingMonad {
	pMonad.ready(modifiedTerm)
	pMonad.ready(postExpression)
	return pMonad
}

func postExpression(pMonad parsingMonad) parsingMonad {
	if pMonad.has('|') {
		pMonad.accept('|')
		pMonad.ready(expression)
	} else if !pMonad.isEmpty() && !pMonad.has('$') {
		pMonad.ready(modifiedTerm)
		pMonad.ready(postExpression)
	}

	return pMonad
}

func modifiedTerm(pMonad parsingMonad) parsingMonad {
	if pMonad.has('(') {
		pMonad.accept('(')
		pMonad.ready(expression)
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept(')')
			return pm
		})

		return pMonad
	}

	pMonad.ready(term)
	pMonad.ready(modifier)
	return pMonad
}

func term(pMonad parsingMonad) parsingMonad {
	if pMonad.has('\\') {
		pMonad.accept('\\')
		pMonad.acceptByte()
	} else if pMonad.has('[') {
		pMonad.ready(class)
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept(']')
			return pm
		})
	} else if pMonad.has('.') {
		pMonad.accept('.')
	} else {
		pMonad.ready(literal)
	}

	return pMonad
}

// optionally added
func modifier(pMonad parsingMonad) parsingMonad {
	if pMonad.has('{') {
		pMonad.ready(rangeModifier)
	} else if pMonad.has('+') {
		pMonad.accept('+')
	} else if pMonad.has('*') {
		pMonad.accept('*')
	}

	return pMonad
}

func rangeModifier(pMonad parsingMonad) parsingMonad {
	pMonad.accept('{')
	if pMonad.has(',') {
		pMonad.accept(',')
		pMonad.ready(decimal)
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept('}')
			return pm
		})

		return pMonad
	}

	pMonad.ready(decimal)
	pMonad.ready(func(pm parsingMonad) parsingMonad {
		pm.accept(',')
		if pm.has('}') {
			pm.accept('}')
		} else {
			pm.ready(decimal)
			pm.ready(func(pm_ parsingMonad) parsingMonad {
				pm_.accept('}')
				return pm_
			})
		}

		return pm
	})

	return pMonad
}

func class(pMonad parsingMonad) parsingMonad {
	for !pMonad.has(']') {
		pMonad.acceptUnicode()
	}

	return pMonad
}

func decimal(pMonad parsingMonad) parsingMonad {
	pMonad.acceptWithin('1', '9')

	for pMonad.within('0', '9') {
		pMonad.acceptWithin('1', '9')
	}

	return pMonad
}

func literal(pMonad parsingMonad) parsingMonad {
	pMonad.acceptUnicode()

	return pMonad
}
