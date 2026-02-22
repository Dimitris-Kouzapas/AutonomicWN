package util

type HashSet[T comparable] struct {
    data map[T]struct{}
}

// NewHashSet creates a new empty HashSet.
func NewHashSet[T comparable]() *HashSet[T] {
    return &HashSet[T]{
        data: make(map[T]struct{}),
    }
}

// Add inserts an element into the HashSet.
func (s *HashSet[T]) Add(value T) {
    s.data[value] = struct{}{}
}

// Remove deletes an element from the HashSet.
func (s *HashSet[T]) Remove(value T) {
    delete(s.data, value)
}

// Contains checks if an element exists in the HashSet.
func (s *HashSet[T]) Contains(value T) bool {
    _, exists := s.data[value]
    return exists
}

// Size returns the number of elements in the HashSet.
func (s *HashSet[T]) Size() int {
    return len(s.data)
}

// Clear removes all elements from the HashSet.
func (s *HashSet[T]) Clear() {
    s.data = make(map[T]struct{})
}

// Slice returns a slice of all elements in the HashSet.
func (s *HashSet[T]) Slice() []T {
    keys := make([]T, 0, len(s.data))
    for key := range s.data {
        keys = append(keys, key)
    }
    return keys
}
