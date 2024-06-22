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

func parseRegex(pMonad parsingMonad) parsingMonad {
	pMonad.acceptIfHas('^')
	pMonad.ready(func(pm parsingMonad) parsingMonad {
		pm.acceptIfHas('$')
		return pm
	})
	pMonad.ready(expression)

	return pMonad
}

func expression(pMonad parsingMonad) parsingMonad {
	if pMonad.isEmpty() || pMonad.has('$') {
		return pMonad
	}

	pMonad.ready(postExpression)
	pMonad.ready(modifiedTerm)
	return pMonad
}

func postExpression(pMonad parsingMonad) parsingMonad {
	if pMonad.has('|') {
		pMonad.accept('|')
		pMonad.ready(expression)
	} else if !pMonad.isEmpty() && !pMonad.has('$') {
		pMonad.ready(postExpression)
		pMonad.ready(modifiedTerm)
	}

	return pMonad
}

func modifiedTerm(pMonad parsingMonad) parsingMonad {
	if pMonad.has('(') {
		pMonad.accept('(')
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept(')')
			return pm
		})
		pMonad.ready(expression)

		return pMonad
	}

	pMonad.ready(modifier)
	pMonad.ready(term)
	return pMonad
}

func term(pMonad parsingMonad) parsingMonad {
	if pMonad.has('\\') {
		pMonad.accept('\\')
		pMonad.acceptByte()
	} else if pMonad.has('[') {
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept(']')
			return pm
		})
		pMonad.ready(class)

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
		pMonad.ready(func(pm parsingMonad) parsingMonad {
			pm.accept('}')
			return pm
		})
		pMonad.ready(decimal)

		return pMonad
	}

	pMonad.ready(func(pm parsingMonad) parsingMonad {
		pm.accept(',')
		if pm.has('}') {
			pm.accept('}')
		} else {
			pm.ready(func(pm_ parsingMonad) parsingMonad {
				pm_.accept('}')
				return pm_
			})
			pm.ready(decimal)
		}

		return pm
	})
	pMonad.ready(decimal)

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
