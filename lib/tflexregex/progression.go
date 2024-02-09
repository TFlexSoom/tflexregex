package tflexregex

type CharRangeProgression struct {
	binaryTree []map[uint]bool
}

const defaultRowCapacity uint = 8

func New() CharRangeProgression {
	return CharRangeProgression{
		binaryTree: make([]map[uint]bool, 512),
	}
}

func (crp *CharRangeProgression) TransitionOnCharacter(charVal byte, item uint) {
	index := uint16(charVal) + 256
	if (*crp).binaryTree[index] == nil {
		(*crp).binaryTree[index] = make(map[uint]bool, defaultRowCapacity)
	}

	(*crp).binaryTree[index][item] = true
}

func (crp *CharRangeProgression) TransitionOnRange(greaterToOrEqual byte, lessThan byte, item uint) {
	if greaterToOrEqual >= lessThan {
		return
	} else if greaterToOrEqual == lessThan-1 {
		crp.TransitionOnCharacter(greaterToOrEqual, item)
		return
	}

	rangeRoot := uint(0)
	index := byte(0b1000_0000)
	for ; index > 0; index = index >> 1 {
		left := greaterToOrEqual & index
		right := lessThan & index

		if left != right {
			break
		} else if left != 0 {
			rangeRoot += uint(left)
		}
	}

	// so now we know greaterToOrEqual & index == 0
	// and lessThan & index == 1
	fullFromIndex := (index << 1) - 1
	left := greaterToOrEqual & fullFromIndex
	right := lessThan & fullFromIndex
}

func (crp *CharRangeProgression) GetTransitions(actual byte) []uint {
	result := make([]uint, 0, 32)
	if crp.binaryTree[0] != nil {
		for k, _ := range crp.binaryTree[0] {
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

		if crp.binaryTree[index] != nil {
			for k, _ := range crp.binaryTree[0] {
				result = append(result, k)
			}
		}
	}

	return result
}
