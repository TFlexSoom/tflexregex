package parse

type Regex struct{}

func parse(pattern string) Regex {
	monad := fromString(pattern)

	monad = runUntilDone(monad)

	return Regex{}
}
