package tflexregex

import "errors"

type RegexBuilder ast

const (
	parseTypeExpression byte = iota
	parseTypeClass
	parseTypeCount
)

type parsingMonad struct {
	astStack     []ast
	parseType    byte
	lastAddition uint
	isEscaped    bool
	isUnioning   bool
	err          error
}

func newParsingMonad() parsingMonad {
	result := parsingMonad{
		astStack:     make([]ast, 0, 8),
		parseType:    parseTypeExpression,
		lastAddition: 0,
		isEscaped:    false,
		isUnioning:   false,
		err:          nil,
	}

	result.astStack = append(result.astStack, newAst())
	return result
}

func parse(regex string) (RegexBuilder, error) {
	pMonad := newParsingMonad()

	for _, v := range regex {
		switch pMonad.parseType {
		case parseTypeExpression:
			pMonad = parseExpression(v, pMonad)
			break
		case parseTypeClass:
			pMonad = parseClass(v, pMonad)
			break
		case parseTypeCount:
			pMonad = parseCount(v, pMonad)
			break
		default:
			return RegexBuilder{}, errors.New("unknown parse type")
		}

		if pMonad.err != nil {
			return RegexBuilder{}, pMonad.err
		}
	}

	if len(pMonad.astStack) > 1 {
		return RegexBuilder{}, errors.New("missing ')' in expression")
	}

	return RegexBuilder(popStack(&pMonad.astStack)), nil
}

func parseExpression(v rune, pMonad parsingMonad) parsingMonad {
	end := len(pMonad.astStack) - 1
	last := &(pMonad.astStack[end])

	switch v {
	case '(':
		(pMonad.astStack) = append(pMonad.astStack, newAst())
		break
	case ')':
		if end == 0 {
			pMonad.err = errors.New("illegal ')' in expression")
			break
		}

		parent := &(pMonad.astStack[end-1])

		if !pMonad.isUnioning {
			pMonad.lastAddition = addSubTree((*parent), popStack(&pMonad.astStack))
		} else {
			pMonad.isUnioning = false
		}

		break
	case '*':
	case '+':
		if len((*last).elem) == 0 {
			pMonad.err = errors.New("cannot modify empty expression")
			break
		}

		lastRoot := &((*last).elem)[pMonad.lastAddition]
		if (*lastRoot).gte != 1 || (*lastRoot).lte != 1 {
			pMonad.err = errors.New("elem already has modifier")
			break
		}

		if v == '*' {
			(*lastRoot).lte = 0
		}

		(*lastRoot).gte = 0
		break
	case '{':
		pMonad.parseType = parseTypeCount
		break
	case '|':
		if len((*last).elem) == 0 {
			pMonad.err = errors.New("illegal '|' in expression")
			break
		}

		return workingCodeUnion, nil
	case '[':
		(*last).classes = append((*last).classes, newClass())
		(*last).series = append((*last).series, newClassElem(uint(len((*last).classes)-1)))
		return workingCodeClass, nil
	case '\\':
		return workingCodeEscape, nil
	case '.':
		(*last).series = append((*last).series, newDotElem())
		break
	default:
		(*last).literals = append((*last).literals, v)
		(*last).series = append((*last).series, newElem(uint(len((*last).literals)-1)))
		break
	}

	return workingCodeExpression, nil
}

func parseClass(v rune, workingstack *[]ast, isEscaped bool) (int, error) {
	if !isEscaped && v == ']' {
		return workingCodeExpression
	}

	series := &((*workingstack)[len(*workingstack)-1].series)

	switch v {
	case '(':
		(*workingstack) = append(*workingstack, newAst())
		break
	case ')':
		if end == 0 {
			return 0, errors.New("illegal ')' in expression")
		}

		parent := &((*workingstack)[end-1])
		(*parent).subs = append((*parent).subs, popStack(workingstack))
		break
	case '*':
	case '+':
		if len((*last).series) == 0 {
			return 0, errors.New("illegal '*' in expression")
		}

		seriesLast := &((*last).series)[len((*last).series)-1]
		if (*seriesLast).secondOrMod != modifierNone {
			return 0, errors.New("elem already has modifier")
		}

		if v == '*' {
			(*seriesLast).secondOrMod = modifierStar
		} else {
			(*seriesLast).secondOrMod = modifierPlus
		}

		break
	case '{':
		return workingCodeCount, nil
	case '|':
		return workingCodeUnion, nil
	case '[':
		return workingCodeClass, nil
	case '\\':
		return workingCodeEscape, nil
	case '.':
		series := &((*last).series)
		(*series) = append((*series), newDotElem())
		break
	default:
		series := &((*last).series)
		(*series) = append((*series), newElem(v))
		break
	}

	return workingCodeExpression, nil
}
