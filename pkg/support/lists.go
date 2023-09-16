package support

func Reverse[T any](list []T) []T {
	// Get the length of the list
	n := len(list)

	// Create a new list of the same type and length as the original list
	reversed := make([]T, n)

	// Iterate over the elements of the original list in reverse order and copy them to the new list
	for i, j := n-1, 0; i >= 0; i, j = i-1, j+1 {
		reversed[j] = list[i]
	}

	return reversed
}
