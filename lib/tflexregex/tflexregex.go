package tflexregex

import (
	"github.com/tflexsoom/tflexregex/lib/tflexregex/parse"
	"github.com/tflexsoom/tflexregex/lib/tflexregex/regex"
)

func Matches(pattern string, b []byte) (bool, error) {
	regex, err := RegexFromString(pattern)
	if err != nil {
		return false, err
	}

	return regex.Matches(b), nil
}

func RegexFromString(pattern string) (regex.Regex, error) {
	v := parse.NewIrVisitor()
	err := parse.Parse(pattern, v)
	if err != nil {
		return nil, err
	}

	return regex.RecursiveTreeFromSingle(v.Ir()), nil
}
