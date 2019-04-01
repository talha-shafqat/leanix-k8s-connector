package main

// StringSet is a helper type to respresent a set of strings
type StringSet struct {
	Map map[string]bool
}

// NewStringSet new StringSet with an empty set
func NewStringSet() *StringSet {
	return &StringSet{
		Map: make(map[string]bool),
	}
}

// Add adds a string to the set
func (s *StringSet) Add(i string) {
	s.Map[i] = true
}

// Items returns all items in the set as slice
func (s *StringSet) Items() []string {
	slice := make([]string, len(s.Map))
	i := 0
	for k := range s.Map {
		slice[i] = k
		i++
	}
	return slice
}
