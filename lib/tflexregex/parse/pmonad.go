package parse

type parsingMonad interface {
	isEmpty() bool
	has(byte) bool
	within(byte, byte) bool
	acceptIfHas(byte)
	accept(byte)
	acceptByte()
	acceptWithin(byte, byte)
	acceptUnicode()
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
