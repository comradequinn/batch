package batch

// stringSet represents a unique list of strings
type stringSet struct {
	m     map[string]struct{}
	Ready chan struct{}
}

func newStringSet() *stringSet {
	return &stringSet{
		m:     map[string]struct{}{},
		Ready: make(chan struct{}, 1),
	}
}

func (s *stringSet) has(element string) bool {
	_, exists := s.m[element]

	return exists
}

func (s *stringSet) add(element string) {
	s.m[element] = struct{}{}
}

func (s *stringSet) size() int {
	return len(s.m)
}
