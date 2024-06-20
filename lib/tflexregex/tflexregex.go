package tflexregex

type Regex struct {
	definition string
}

// type regexMatchMonad struct {
// 	regex Regex
// 	state uint
// }

// func fromRegex(regex Regex) regexMonad {
// 	return regexMonad{
// 		regex:  regex,
// 		founds: make([]uint, 0, 16),
// 	}
// }

// func hasFound(monad regexMonad) bool {
// 	return len(monad.founds) > 0
// }

// func match(monad regexMonad, next byte) regexMonad {

// }

// func monadMatch(regex Regex, b []byte) bool {
// 	monad := fromRegex(regex)
// 	for i := 0; i < len(b); i++ {
// 		monad = match(monad, b[i])
// 	}

// 	return hasFound(monad)
// }

func Matches(pattern string, b []byte) (bool, error) {
	return false, nil
}

func (regex *Regex) Matches(b []byte) bool {
	return false
}

func RegexFromString(pattern string) (Regex, error) {
	regex := NewRegex()
	// TODO fix
	regex.definition = pattern
	return regex, nil
}

func NewRegex() Regex {
	return Regex{}
}
