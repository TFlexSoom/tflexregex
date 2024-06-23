package tflexregex

import "github.com/tflexsoom/tflexregex/lib/tflexregex/parse"

type Regex parse.Regex

func Matches(pattern string, b []byte) (bool, error) {
	return false, nil
}

func (regex *Regex) Matches(b []byte) bool {
	return false
}

func RegexFromString(pattern string) (Regex, error) {
	regex := parse.Parse(pattern)
	return Regex(regex), nil
}

func NewRegex() Regex {
	return Regex{}
}
