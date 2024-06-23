package parse

const MAX_CALLSTACK = 256
const DEFAULT_CALLSTACK_SIZE = 32

type parsingMonadString struct {
	content   []byte
	length    int
	index     int
	callstack [](func(parsingMonad) parsingMonad)
	visit     visitor
}

func fromString(pattern string) parsingMonad {
	bytes := []byte(pattern)

	return &parsingMonadString{
		content:   bytes,
		length:    len(bytes),
		index:     0,
		callstack: make([]func(parsingMonad) parsingMonad, 0, 32),
	}
}

func (pm parsingMonadString) isEmpty() bool {
	return pm.index < pm.length
}

func (pm parsingMonadString) has(b byte) bool {
	if pm.isEmpty() {
		panic("cannot check empty monad")
	}

	return pm.content[pm.index] == b
}

func (pm parsingMonadString) within(begin byte, end byte) bool {
	if pm.isEmpty() {
		panic("cannot check empty monad")
	}

	current := pm.content[pm.index]

	return begin <= current && end >= current
}

func (pm *parsingMonadString) accept(b byte) {
	if pm.isEmpty() {
		panic("cannot accept on empty monad")
	}

	if pm.has(b) {
		panic("byte does not match expected byte")
	}

	pm.acceptByte()
}

func (pm *parsingMonadString) acceptIfHas(b byte) {
	if pm.isEmpty() {
		panic("cannot accept on empty monad")
	}

	if pm.has(b) {
		return
	}

	pm.acceptByte()
}

func (pm *parsingMonadString) acceptByte() {
	if pm.isEmpty() {
		panic("cannot accept on empty monad")
	}

	(*pm).index += 1
}

func (pm *parsingMonadString) acceptWithin(begin byte, end byte) {
	if pm.isEmpty() {
		panic("cannot accept on empty monad")
	}

	if pm.within(begin, end) {
		panic("provided character is not within bounds")
	}

	pm.acceptByte()
}

func (pm *parsingMonadString) acceptUnicode() {
	for i := 0; i < 4; i++ {
		if pm.isEmpty() {
			panic("cannot accept on empty monad")
		}

		if pm.within(byte(0), byte(127)) {
			pm.acceptByte()
			return
		}

		pm.acceptByte()
	}

	panic("bad unicode character")
}

func (pm *parsingMonadString) pump(nt nodeType) {
	(*pm).visit.pump(nt)
}

func (pm *parsingMonadString) ready(nextTerm func(parsingMonad) parsingMonad) {
	pmonad := *pm
	if len(pmonad.callstack) > MAX_CALLSTACK {
		panic("recursive descent went past MAX_CALLSTACK")
	}

	pm.callstack = append(pm.callstack, nextTerm)
	*pm = pmonad
}

func (pm parsingMonadString) isRunning() bool {
	return len(pm.callstack) > 0
}

func (pm *parsingMonadString) nextFunction() func(parsingMonad) parsingMonad {
	pMonad := *pm
	length := len(pMonad.callstack)
	result := pMonad.callstack[length-1]
	pMonad.callstack = pMonad.callstack[:length-1]
	*pm = pMonad
	return result
}
