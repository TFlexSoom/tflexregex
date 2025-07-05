package parse

type Regex struct{}

func Parse(pattern string, v visitor) error {
	monad := from(pattern, v)
	return runUntilDone(monad)
}
