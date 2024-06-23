package parse

type Regex struct{}

func Parse(pattern string) Regex {
	monad := fromString(pattern)

	monad = runUntilDone(monad)

	return Regex{}
}
