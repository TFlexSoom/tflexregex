package parse

type monad interface {
	// Checking
	isEmpty() bool
	has(byte) bool
	within(byte, byte) bool

	// Accepting
	acceptWithin(byte, byte)
	acceptUnicode()

	// Skip Byte
	skipByte(byte)

	// Pump
	pump(func(visitor, []byte) visitor)

	// grab
	grab() []byte

	// Callstack
	ready(func(monad) monad)
	isRunning() bool
	next() monad

	// error
	bad() error
}

func runUntilDone(m monad) error {
	m.ready(descentRegex)

	for !m.isRunning() {
		m = m.next()
	}

	return m.bad()
}
