package tflexregex

import (
	"log"
	"strconv"
)

const (
	RANGE_STATE_FIND_LOWER_BEG = iota
	RANGE_STATE_FIND_LOWER_END
	RANGE_STATE_FIND_COMMA
	RANGE_STATE_FIND_UPPER_BEG
	RANGE_STATE_FIND_UPPER_END
	RANGE_STATE_FIND_END_RANGE
	END_STATE
)

type recursiveRange struct {
	lower uint16
	upper uint16
	flags uint8
}

func zeroRange() recursiveRange {
	return recursiveRange{
		lower: 0,
		upper: 0,
		flags: 0,
	}
}

func basicRange(lower uint16, upper uint16) recursiveRange {
	return recursiveRange{
		lower: lower,
		upper: upper,
		flags: 0,
	}
}

func lowerToInfi(lower uint16) recursiveRange {
	return recursiveRange{
		lower: lower,
		upper: 0,
		flags: 1,
	}
}

func capFromInfi(upper uint16) recursiveRange {
	return recursiveRange{
		lower: 0,
		upper: upper,
		flags: 2,
	}
}

func getRange(regex string, regexIndex uint) recursiveRange {
	regexLen := uint(len(regex))
	lowerBeg := uint(0)
	lowerEnd := uint(0)
	upperBeg := uint(0)
	upperEnd := uint(0)
	state := RANGE_STATE_FIND_LOWER_BEG

	for i := regexIndex; i < regexLen && state <= RANGE_STATE_FIND_END_RANGE; i++ {
		switch state {
		case RANGE_STATE_FIND_LOWER_BEG:
			if regex[i] == ' ' || regex[i] == '{' {
				continue
			} else if regex[i] == ',' {
				state = RANGE_STATE_FIND_UPPER_BEG
			} else {
				lowerBeg = i
				state = RANGE_STATE_FIND_LOWER_END
			}
			break
		case RANGE_STATE_FIND_LOWER_END:
			if regex[i] == ',' {
				lowerEnd = i
				state = RANGE_STATE_FIND_UPPER_BEG
			}
			if regex[i] < '0' || regex[i] > '9' {
				lowerEnd = i
				state = RANGE_STATE_FIND_COMMA
			}
			break
		case RANGE_STATE_FIND_COMMA:
			if regex[i] == ',' {
				state = RANGE_STATE_FIND_UPPER_BEG
			}
			break
		case RANGE_STATE_FIND_UPPER_BEG:
			if regex[i] == ' ' || regex[i] == '{' {
				continue
			} else if regex[i] == '}' {
				state = END_STATE
			} else {
				upperBeg = i
				state = RANGE_STATE_FIND_LOWER_END
			}
			break
		case RANGE_STATE_FIND_UPPER_END:
			if regex[i] == '}' {
				upperEnd = i
				state = END_STATE
			}
			if regex[i] < '0' || regex[i] > '9' {
				lowerEnd = i
				state = RANGE_STATE_FIND_END_RANGE
			}
			break
		case RANGE_STATE_FIND_END_RANGE:
			if regex[i] == '}' {
				state = END_STATE
			}
			break
		default:
			state = END_STATE
		}
	}

	isInfiLower := lowerEnd-lowerBeg == 0
	isInfiUpper := upperEnd-upperBeg == 0

	if isInfiLower && isInfiUpper {
		return zeroRange()
	} else if isInfiLower {
		upperStr := string(regex[upperBeg:upperEnd])
		upper, err := strconv.Atoi(upperStr)
		if err != nil {
			log.Printf("Error from converting decimal %v: %v", upperStr, err)
		}
		return capFromInfi(uint16(upper))
	} else if isInfiUpper {
		lowerStr := string(regex[lowerBeg:lowerEnd])
		lower, err := strconv.Atoi(lowerStr)
		if err != nil {
			log.Printf("Error from converting decimal %v: %v", lowerStr, err)
		}
		return lowerToInfi(uint16(lower))
	} else {
		upperStr := string(regex[upperBeg:upperEnd])
		upper, err := strconv.Atoi(upperStr)
		if err != nil {
			log.Printf("Error from converting decimal %v: %v", upperStr, err)
		}

		lowerStr := string(regex[lowerBeg:lowerEnd])
		lower, err := strconv.Atoi(lowerStr)
		if err != nil {
			log.Printf("Error from converting decimal %v: %v", lowerStr, err)
		}

		return basicRange(uint16(lower), uint16(upper))
	}

}
