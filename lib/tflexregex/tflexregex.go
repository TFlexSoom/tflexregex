package tflexregex

import (
	"errors"
)

type Regex struct {
}

func FromString(regex string) (Regex, error) {
	builder, err := builderFrom(regex)
	if err != nil {
		return Regex{}, err
	}

	return FromBuilder(builder)
}

func FromBuilder(regex RegexBuilder) (Regex, error) {
	return Regex{}, errors.New("not implemented")
}

func builderFrom(regex string) (RegexBuilder, error) {
	return RegexBuilder{}, errors.New("not implemented")
}
