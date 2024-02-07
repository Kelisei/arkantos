package internal

type State struct {
	Buffers        []Buffer
	CurrentBuffer  int
	ProjectPath    string
	lastUpdateTime int64
}

func (s *State) cycleBuffer() {
	if s.CurrentBuffer < len(s.Buffers)-1 {
		s.CurrentBuffer++
	} else {
		s.CurrentBuffer = 0
	}
}
