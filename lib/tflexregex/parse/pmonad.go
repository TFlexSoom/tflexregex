package parse

type nodeType byte

const (
	NODE_TYPE_NONE nodeType = iota
	NODE_TYPE_ANCHOR
	NODE_TYPE_DOLLAR

	NODE_TYPE_CHAR
	NODE_TYPE_UNICODE
	NODE_TYPE_CLASS
	NODE_TYPE_DOT

	NODE_TYPE_MODIFIER

	NODE_TYPE_PAR_BEGIN
	NODE_TYPE_PAR_END

	NODE_TYPE_UNION
)

type parsingMonad interface {
	// Checking
	isEmpty() bool
	has(byte) bool
	within(byte, byte) bool

	// Accepting
	acceptIfHas(byte)
	accept(byte)
	acceptByte()
	acceptWithin(byte, byte)
	acceptUnicode()

	// Buffer
	pump(nodeType)
	pumpIfAccepted(nodeType)

	// Callstack
	ready(func(parsingMonad) parsingMonad)
	isRunning() bool
	nextFunction() func(parsingMonad) parsingMonad
}

func runUntilDone(monad parsingMonad) parsingMonad {
	monad.ready(parseRegex)

	for !monad.isRunning() {
		monad = monad.nextFunction()(monad)
	}

	return monad
}

// TODO could add a runStep function
