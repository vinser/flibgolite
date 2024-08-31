package fb2

import "sync"

type tagStack struct {
	mx  sync.Mutex
	tag []string
}

func newTagStack() *tagStack {
	return &tagStack{sync.Mutex{}, make([]string, 0)}
}
func (s *tagStack) push(name string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.tag = append(s.tag, name)
}

func (s *tagStack) pop() string {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(s.tag) == 0 {
		return ""
	}
	item := s.tag[len(s.tag)-1]
	s.tag = s.tag[:len(s.tag)-1]
	return item
}

func (s *tagStack) top() string {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(s.tag) == 0 {
		return ""
	}
	return s.tag[len(s.tag)-1]
}

func (s *tagStack) reset() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.tag = make([]string, 0)
}
