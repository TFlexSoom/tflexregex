package tflexregex

import (
	"fmt"
)

type Progression struct {
	binaryTreeOfSets []map[uint]bool
}

const defaultRowCapacity uint = 8

func NewProgression() Progression {
	return Progression{
		binaryTreeOfSets: make([]map[uint]bool, 511),
	}
}

func (crp *Progression) setTransition(index uint, nextState uint) {
	if (*crp).binaryTreeOfSets[index] == nil {
		(*crp).binaryTreeOfSets[index] = make(map[uint]bool, defaultRowCapacity)
	}

	(*crp).binaryTreeOfSets[index][nextState] = true
}

func (crp *Progression) TransitionOnCharacter(charVal byte, nextState uint) {
	crp.setTransition(uint(charVal)+255, nextState)
}

func (crp *Progression) TransitionOnRange(greaterToOrEqual byte, lessThanOrEqual byte, nextState uint) error {
	if greaterToOrEqual > lessThanOrEqual {
		return fmt.Errorf("invalid range: %d must be less than or equal to %d", greaterToOrEqual, lessThanOrEqual)
	} else if greaterToOrEqual == lessThanOrEqual {
		crp.TransitionOnCharacter(greaterToOrEqual, nextState)
		return nil
	} else if greaterToOrEqual == 0 && lessThanOrEqual == 255 {
		crp.setTransition(0, nextState)
	}

	upperBits := byte(0)
	iter := byte(0b1000_0000)
	for ; iter > 0; iter = iter >> 1 {
		left := greaterToOrEqual & iter
		right := lessThanOrEqual & iter

		if left != right {
			break
		} else {
			upperBits += left
		}
	}

	lowerBits := (iter << 1) - 1
	queue := make([]byte, 0, 256)
	queue = append(queue, lowerBits)
	for ; len(queue) > 0; iter = queue[0] {
		lowest := upperBits
		highest := upperBits + lowerBits
		if greaterToOrEqual <= lowest && lessThanOrEqual >= highest {
			crp.setTransition( /* TODO */ 0, nextState)
		} else if greaterToOrEqual <= lowest {
			queue = append(queue, index) // TODO
		} else if lessThanOrEqual >= highest {
			queue = append(queue, index) // TODO
		}
	}
}

func (crp *Progression) GetTransitions(actual byte) []uint {
	result := make([]uint, 0, 32)
	if crp.binaryTreeOfSets[0] != nil {
		for k, _ := range crp.binaryTreeOfSets[0] {
			result = append(result, k)
		}
	}

	index := uint(0)
	for bitwise := byte(0b1000_0000); bitwise > 0; bitwise = bitwise >> 1 {
		if actual&bitwise != 0 {
			index = 2*index + 2
		} else {
			index = 2*index + 1
		}

		if crp.binaryTreeOfSets[index] != nil {
			for k, _ := range crp.binaryTreeOfSets[0] {
				result = append(result, k)
			}
		}
	}

	return result
}
