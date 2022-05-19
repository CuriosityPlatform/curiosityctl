package slices

// Map iterates through slice and maps values
func Map[T, TResult any](s []T, f func(T) TResult) []TResult {
	result := make([]TResult, 0, len(s))
	for _, t := range s {
		result = append(result, f(t))
	}
	return result
}

// MapErr iterates through slice and maps values and stops on any error
func MapErr[T, TResult any](s []T, f func(T) (TResult, error)) ([]TResult, error) {
	result := make([]TResult, 0, len(s))
	for _, t := range s {
		e, err := f(t)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// Filter iterates and adds to result slice elements that satisfied predicate
func Filter[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, t := range s {
		if f(t) {
			result = append(result, t)
		}
	}
	return result
}

// FilterErr iterates and adds to result slice elements that satisfied predicate and stop on any error
func FilterErr[T any](s []T, f func(T) (bool, error)) ([]T, error) {
	var result []T
	for _, t := range s {
		accepted, err := f(t)
		if err != nil {
			return nil, err
		}

		if accepted {
			result = append(result, t)
		}
	}

	return result, nil
}

// Chunk chunks slice of items into slices by passed size
func Chunk[T any](items []T, size int) (chunks [][]T) {
	i := 0
	counter := 0
	chunk := make([]T, 0, size)
	for _, item := range items {
		i++
		counter++

		chunk = append(chunk, item)

		if i == size {
			chunks = append(chunks, chunk)
			i = 0

			// Allocate only estimated count of elements
			if estimated := len(items) - counter; (estimated) < size {
				chunk = make([]T, 0, estimated)
				continue
			}
			chunk = make([]T, 0, size)
		}
	}

	if len(chunk) != 0 {
		chunks = append(chunks, chunk)
	}

	return
}
