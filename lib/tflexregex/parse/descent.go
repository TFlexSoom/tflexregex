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
LITERAL => [^()'\'+*{}[\].]
DECIMAL => 123456789+
DOT => '.'
ESCAPE_LITERAL => '\'
*/

// REGEX => '^'? EXPRESSION? '$'?
func descentRegex(pMonad monad) monad {
	if pMonad.has('^') {
		pMonad.skipByte('^')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.anchor()
		})
	}

	pMonad.ready(optionalDollar)
	pMonad.ready(optionalExpression)

	return pMonad
}

func optionalExpression(pMonad monad) monad {
	if pMonad.isEmpty() || pMonad.has('$') {
		return pMonad
	}

	pMonad.ready(expression)
	return pMonad
}

func optionalDollar(pMonad monad) monad {
	if pMonad.has('$') {
		pMonad.skipByte('$')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.dollar()
		})
	}
	return pMonad
}

// EXPRESSION => MODIFIED_TERM + ('|' EXPRESSION )*
func expression(pMonad monad) monad {
	if pMonad.isEmpty() {
		panic("parsing error no modified terms")
	}

	pMonad.ready(starExpression)
	pMonad.ready(modifiedTerm)
	return pMonad
}

func starExpression(pMonad monad) monad {
	if pMonad.has('|') {
		pMonad.skipByte('|')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.union()
		})
		pMonad.ready(expression)
	} else if !pMonad.isEmpty() && !pMonad.has('$') {
		pMonad.ready(starExpression)
		pMonad.ready(modifiedTerm)
	}

	return pMonad
}

// MODIFIED_TERM => ('(' EXPRESSION ')' MODIFIER?) | (TERM MODIFIER?)
func modifiedTerm(pMonad monad) monad {
	if pMonad.has('(') {
		pMonad.skipByte('(')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.openParenthesis()
		})
		pMonad.ready(func(pm monad) monad {
			pm.skipByte(')')
			pm.pump(func(v visitor, bs []byte) visitor {
				return v.closeParenthesis()
			})
			return pm
		})
		pMonad.ready(expression)

		return pMonad
	}

	pMonad.ready(modifier)
	pMonad.ready(term)
	return pMonad
}

func optionalModifier(pMonad monad) monad {
	if !pMonad.has('{') && !pMonad.has('+') && !pMonad.has('*') {
		return pMonad
	}

	pMonad.ready(modifier)
	return pMonad
}

// TERM => ESCAPE | CLASS | DOT | LITERAL
func term(pMonad monad) monad {
	if pMonad.has('\\') {
		pMonad.skipByte('\\')
		pMonad.acceptWithin(byte(0), byte(255))
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.char(bs[0])
		})
	} else if pMonad.has('[') {
		pMonad.ready(func(pm monad) monad {
			pm.skipByte(']')
			return pm
		})
		pMonad.ready(class)

	} else if pMonad.has('.') {
		pMonad.skipByte('.')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.dot()
		})
	} else {
		pMonad.ready(literal)
	}

	return pMonad
}

// MODIFIER => {DECIMAL,DECIMAL} | {DECIMAL,} | {,DECIMAL} | '+' | '*'
func modifier(pMonad monad) monad {
	if pMonad.has('{') {
		pMonad.ready(rangeModifier)
		return pMonad
	} else if pMonad.has('+') {
		pMonad.skipByte('+')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.modifier(1, 0)
		})
		return pMonad
	} else if pMonad.has('*') {
		pMonad.skipByte('*')
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.modifier(0, 0)
		})
		return pMonad
	}

	panic("modifier expected but not found")
}

func rangeModifier(pMonad monad) monad {
	pMonad.skipByte('{')
	if pMonad.has(',') {
		pMonad.skipByte(',')
		pMonad.ready(func(pm monad) monad {
			pMonad.pump(func(v visitor, bs []byte) visitor {
				return v.modifier(0, uint(atoi(bs)))
			})
			pm.skipByte('}')
			return pm
		})
		pMonad.ready(decimal)

		return pMonad
	}

	pMonad.ready(func(pm monad) monad {
		pm.skipByte(',')
		if pm.has('}') {
			pm.pump(func(v visitor, bs []byte) visitor {
				return v.modifier(uint(atoi(bs)), 0)
			})
			pm.skipByte('}')
		} else {
			lower := uint(atoi(pm.grab()))
			pm.ready(func(pm_ monad) monad {
				pm_.pump(func(v visitor, bs []byte) visitor {
					return v.modifier(lower, uint(atoi(bs)))
				})
				pm_.skipByte('}')
				return pm_
			})
			pm.ready(decimal)
		}

		return pm
	})
	pMonad.ready(decimal)

	return pMonad
}

// CLASS => '[' [^]]* ']'
func class(pMonad monad) monad {
	pMonad.skipByte('[')
	for !pMonad.has(']') {
		pMonad.acceptUnicode()
		pMonad.pump(func(v visitor, bs []byte) visitor {
			return v.unicode(utf8(bs))
		})
	}

	return pMonad
}

// DECIMAL => 123456789+
func decimal(pMonad monad) monad {
	pMonad.acceptWithin('1', '9')

	for pMonad.within('0', '9') {
		pMonad.acceptWithin('0', '9')
	}

	return pMonad
}

// LITERAL => [^()'\'+*{}[\].]
func literal(pMonad monad) monad {
	if pMonad.has('^') ||
		pMonad.has('(') ||
		pMonad.has(')') ||
		pMonad.has('\\') ||
		pMonad.has('+') ||
		pMonad.has('*') ||
		pMonad.has('{') ||
		pMonad.has('}') ||
		pMonad.has('[') ||
		pMonad.has(']') ||
		pMonad.has('.') {
		panic("bad literal term")
	}

	pMonad.acceptUnicode()
	pMonad.pump(func(v visitor, bs []byte) visitor {
		return v.unicode(utf8(bs))
	})

	return pMonad
}
