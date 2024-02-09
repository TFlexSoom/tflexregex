package tflexregex

import (
	"errors"
)

type Regex struct {
	currentstate CharacterProgression
	states       [][]filter
}

const (
	FILTER_FLAG_GTE    uint8 = 1 << 0
	FILTER_FLAG_GTE_LT uint8 = 1 << 1
)

type filter struct {
	gte    byte
	lt     byte
	states []uint
	flags  uint8
}

///////////////////////////////////////////////////////////////////////

const stateCap = 32
const filtersCap = 32
const literalFilterCap = 32
const classFilterCap = 32
const stateTransitionCap = 8
const callbackCap = 8

type rIndex struct {
	regexIndex uint
	stateIndex uint
}

type readyState struct {
	states     []state
	stateIndex uint
	statesLen  uint

	filters    []filter
	filtersLen uint

	literalFilters map[byte]uint
	classFilters   map[string]uint

	stateStack []uint
	stackIndex uint

	backslashFlag bool
	unicodeFlag   uint

	regex       string
	regexIndex  uint
	regexLength uint

	rCalls      []rIndex
	rCallLength uint
}

func Ready(regex string) (Regex, error) {
	return recurseReady(readyState{
		states:         make([]state, 0, stateCap),
		stateIndex:     0,
		statesLen:      0,
		filters:        make([]filter, 0, 32),
		filtersLen:     0,
		literalFilters: make(map[byte]uint, literalFilterCap),
		classFilters:   make(map[string]uint, classFilterCap),
		stateStack:     make([]uint, 32),
		stackIndex:     0,
		backslashFlag:  false,
		unicodeFlag:    0,
		regex:          regex,
		regexIndex:     0,
		regexLength:    uint(len(regex)),
		rCalls:         make([]rIndex, callbackCap),
		rCallLength:    0,
	})
}

func recurseReady(rState readyState) (Regex, error) {
	rState.rCalls[0] = rIndex{
		regexIndex: 0,
		stateIndex: 0,
	}

	var e error
	for rState.rCallLength = 1; rState.rCallLength > 0; rState.rCallLength-- {
		rState.regexIndex = rState.rCalls[rState.rCallLength-1].regexIndex
		rState.stateIndex = rState.rCalls[rState.rCallLength-1].stateIndex

		rState, e = readySwitch(rState)
		if e != nil {
			return Regex{}, e
		}
	}

	return buildFromState(rState), nil
}

func readySwitch(rState readyState) (readyState, error) {
	for ; rState.regexIndex < rState.regexLength; rState.regexIndex++ {
		v := rState.regex[rState.regexIndex]

		if rState.backslashFlag {
			rState.backslashFlag = false
			rState = getOrAppendLiteralFilter(
				rState,
				v,
				filter{
					gte:   v,
					lt:    v + 1,
					flags: FILTER_FLAG_GTE_LT,
				})
			continue
		}

		switch v {
		case '(':
			break
		case ')':
			break
		case '*':
			break
		case '+':
			break
		case '{':
			rRange := getRange(rState.regex, rState.regexIndex)
			break
		case '|':
			//TODO
			break
		case '.':
			rState = getOrAppendLiteralFilter(
				rState,
				v,
				filter{
					gte:   0,
					flags: FILTER_FLAG_GTE,
				})
			break
		case '[':
			var err error
			rState, err = getOrAppendClassFilter(rState)
			if err != nil {
				return rState, err
			}
			break
		case '\\':
			rState.backslashFlag = true
			break
		default:
			rState = getOrAppendLiteralFilter(
				rState,
				v,
				filter{
					gte:   v,
					lt:    v + 1,
					flags: FILTER_FLAG_GTE,
				})
		}
	}

	return rState, nil
}

func buildFromState(rState readyState) Regex {
	// TODO
	return Regex{}
}

func getOrAppendLiteralFilter(
	rState readyState,
	index byte,
	f filter,
) readyState {
	filterIndex, exists := rState.literalFilters[index]
	if !exists {
		rState.filters = append(rState.filters, filter{
			gte:   0,
			flags: FILTER_FLAG_GTE,
		})

		rState.literalFilters[index] = rState.filtersLen
		filterIndex = rState.filtersLen
		rState.filtersLen += 1
	}

	if rState.stateIndex > rState.statesLen {
		stateTransitions := make([]uint, 1, stateTransitionCap)
		stateTransitions[0] = filterIndex
		rState.states = append(rState.states, state{
			filters: stateTransitions,
		})

		rState.statesLen += 1
	} else {
		rState.states[rState.stateIndex].filters = append(rState.states[rState.stateIndex].filters, filterIndex)
	}

	return rState
}

func getOrAppendClassFilter(
	rState readyState,
) (readyState, error) {
	begin := rState.regexIndex
	utf8Flag := 0
	flag := false

	for rState.regexIndex += 1; rState.regexIndex < rState.regexLength; (rState.regexIndex)++ {
		if utf8Flag > 0 {
			utf8Flag -= 1
			continue
		}

		if rState.regex[rState.regexIndex] == '\\' {
			flag = true
			continue
		} else if rState.regex[rState.regexIndex] == ']' && !flag {
			break
		}

		flag = false
	}

	if rState.regexIndex >= rState.regexLength {
		return rState, errors.New("missing matching ] in character class")
	} else if rState.regexIndex == begin+1 {
		return rState, errors.New("empty character class")
	}

	// TODO
	return rState, nil
}
