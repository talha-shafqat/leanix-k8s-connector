package set

// String is a helper type to represent a set of strings
type String struct {
	Map map[string]bool
}

// NewStringSet new StringSet with an empty set
func NewStringSet() *String {
	return &String{
		Map: make(map[string]bool),
	}
}

// Add adds a string to the set
func (s *String) Add(i string) {
	s.Map[i] = true
}

// Items returns all items in the set as slice
func (s *String) Items() []string {
	slice := make([]string, len(s.Map))
	i := 0
	for k := range s.Map {
		slice[i] = k
		i++
	}
	return slice
}

// Contains returns true if the set contains the given string
func (s *String) Contains(input string) bool {
	_, ok := s.Map[input]
	return ok
}
