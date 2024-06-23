package progression

import "fmt"

type fixedBTree[T any] struct {
	bTree []T
}

func NewFixedBTree[T any]() fixedBTree[T] {
	return fixedBTree[T]{
		bTree: make([]T, 512),
	}
}

func (bt *fixedBTree[T]) mapOn(index uint, mapper func(*T)) {
	mapper(&(*bt).bTree[index])
}

func treeIndexToRange(index uint) (byte, byte) {
	left := index
	for ; left < 255; left = (left << 1) + 1 {
	}
	right := index
	for ; right < 255; right = (right << 1) + 2 {
	}
	return byte(left - 255), byte(right - 255)
}

func (bt *fixedBTree[T]) mapOnLeafNode(index byte, mapper func(*T)) {
	bt.mapOn(uint(index)+255, mapper)
}

func (bt *fixedBTree[T]) mapOnRange(greaterToOrEqual byte, lessThanOrEqual byte, mapper func(*T)) error {
	if greaterToOrEqual > lessThanOrEqual {
		return fmt.Errorf("invalid range: %d must be less than or equal to %d", greaterToOrEqual, lessThanOrEqual)
	} else if greaterToOrEqual == lessThanOrEqual {
		bt.mapOnLeafNode(greaterToOrEqual, mapper)
		return nil
	} else if greaterToOrEqual == 0 && lessThanOrEqual == 255 {
		bt.mapOn(0, mapper)
	}

	treeIndex := uint(0)

	for iter := byte(0b1000_0000); iter > 0; iter = iter >> 1 {
		left := greaterToOrEqual & iter
		right := lessThanOrEqual & iter
		one_zero := uint(0)
		if left > 1 {
			one_zero = 1
		}

		if left != right {
			break
		} else {
			treeIndex = (treeIndex << 1) + one_zero
		}
	}

	queue := make([]uint, 0, 256)
	queue = append(queue, treeIndex)
	for ; len(queue) > 0; treeIndex = popQueue(&queue) {
		lowest, highest := treeIndexToRange(treeIndex)
		if greaterToOrEqual <= lowest && lessThanOrEqual >= highest {
			bt.mapOn(treeIndex, mapper)
		} else if greaterToOrEqual <= lowest {
			queue = append(queue, (treeIndex<<1)+1)
		} else if lessThanOrEqual >= highest {
			queue = append(queue, (treeIndex<<1)+2)
		}
	}

	return nil
}

func foldOnPath[T any, R any, O any](bt *fixedBTree[T], path byte, mapper func(*T) R, fold func(R, O) O, base func() O) O {
	result := base()
	result = fold(mapper(&(*bt).bTree[0]), result)

	index := uint(0)
	for bitwise := byte(0b1000_0000); bitwise > 0; bitwise = bitwise >> 1 {
		if path&bitwise != 0 {
			index = 2*index + 2
		} else {
			index = 2*index + 1
		}

		result = fold(mapper(&(*bt).bTree[index]), result)
	}

	return result
}
