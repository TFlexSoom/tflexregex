package tflexregex

type Regex struct {
}

type regexMatchMonad struct {
	regex Regex
	state uint
}

func fromRegex(regex Regex) regexMonad {
	return regexMonad{
		regex:  regex,
		founds: make([]uint, 0, 16),
	}
}

func hasFound(monad regexMonad) bool {
	return len(monad.founds) > 0
}

func match(monad regexMonad, next byte) regexMonad {

}

// TODO MOVE elsewhere
func Match(regex Regex, b []byte) bool {
	monad := fromRegex(regex)
	for i := 0; i < len(b); i++ {
		monad = match(monad, b[i])
	}

	return hasFound(monad)
}

func NewInstance() Regex {
	return Regex{}
}
