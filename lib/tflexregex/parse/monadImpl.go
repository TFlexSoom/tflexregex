package parse

const MAX_CALLSTACK = 256
const DEFAULT_CALLSTACK_SIZE = 32

type monadImpl struct {
	content   []byte
	length    uint
	left      uint
	right     uint
	visit     visitor
	callstack [](func(monad) monad)
	err       error
}

func from(pattern string, visit visitor) monad {
	bytes := []byte(pattern)

	return &monadImpl{
		content:   bytes,
		length:    uint(len(bytes)),
		left:      0,
		right:     0,
		visit:     visit,
		callstack: make([]func(monad) monad, 0, 32),
		err:       nil,
	}
}

func (m monadImpl) isEmpty() bool {
	return m.right < m.length
}

func (m monadImpl) has(b byte) bool {
	if m.isEmpty() {
		panic("cannot check empty monad")
	}

	return m.content[m.right] == b
}

func (m monadImpl) within(begin byte, end byte) bool {
	if m.isEmpty() {
		panic("cannot check empty monad")
	}

	current := m.content[m.right]

	return begin <= current && end >= current
}

func (m *monadImpl) acceptWithin(begin byte, end byte) {
	if m.isEmpty() {
		panic("cannot accept on empty monad")
	}

	if m.within(begin, end) {
		panic("provided character is not within bounds")
	}

	(*m).right += 1
}

func (m *monadImpl) acceptUnicode() {
	for i := 0; i < 4; i++ {
		if m.isEmpty() {
			panic("cannot accept on empty monad")
		}

		if m.within(byte(0), byte(127)) {
			(*m).right += 1
			return
		}

		(*m).right += 1
	}

	panic("bad unicode character")
}

func (m *monadImpl) skipByte(b byte) {
	if (*m).isEmpty() {
		panic("cannot accept on empty monad")
	}

	if !(*m).has(b) {
		panic("provided character is not expected skip")
	}

	if (*m).left+1 != (*m).right {
		panic("cannot skip with bytes in buffer")
	}

	(*m).right += 1
	(*m).left = (*m).right
}

func (m *monadImpl) pump(termCall func(visitor, []byte) visitor) {
	(*m).visit = termCall((*m).visit, (*m).grab())
}

func (m *monadImpl) grab() []byte {
	if m.left == m.right {
		return []byte{}
	}

	result := m.content[m.left:m.right]

	(*m).left = (*m).right

	return result
}

func (m *monadImpl) ready(nextTerm func(monad) monad) {
	pmonad := *m
	if len(pmonad.callstack) > MAX_CALLSTACK {
		panic("recursive descent went past MAX_CALLSTACK")
	}

	m.callstack = append(m.callstack, nextTerm)
	*m = pmonad
}

func (m monadImpl) isRunning() bool {
	return len(m.callstack) > 0 && m.err == nil
}

func (m *monadImpl) next() monad {
	length := len((*m).callstack)
	m = ((*m).callstack[length-1](m).(*monadImpl))
	(*m).callstack = (*m).callstack[:length-1]
	return m
}

func (m monadImpl) bad() error {
	return m.err
}
