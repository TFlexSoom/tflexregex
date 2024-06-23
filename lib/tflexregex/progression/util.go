package progression

func popQueue[T interface{}](queue *[]T) T {
	result := (*queue)[0]
	*queue = (*queue)[1:]
	return result
}

func popStack[T interface{}](stack *[]T) T {
	last := len(*stack) - 1
	result := (*stack)[last]
	*stack = (*stack)[:last-1]
	return result
}
