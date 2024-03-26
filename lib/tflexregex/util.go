package tflexregex

func popQueue[T interface{}](queue *[]T) T {
	result := (*queue)[0]
	*queue = (*queue)[1 : len(*queue)-1]
	return result
}

func popStack[T interface{}](stack *[]T) T {
	result := (*stack)[len(*stack)-1]
	*stack = (*stack)[1 : len(*stack)-2]
	return result
}
